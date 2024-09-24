package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	_ "github.com/mkdtemplar/simplebank-new/doc/statik"
	"github.com/mkdtemplar/simplebank-new/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/mkdtemplar/simplebank-new/api"
	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
	"github.com/mkdtemplar/simplebank-new/gapi"
	"github.com/mkdtemplar/simplebank-new/pb"
	"github.com/mkdtemplar/simplebank-new/util"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load configuration")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to database")
	}

	//runDbMigrations(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)

	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("gRPC server is listening on %s\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot start server %v", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot start server %v", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("HTTP gateway server is listening on %s\n", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

}

func runDbMigrations(migrationURL string, dbSource string) {
	migrations, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = migrations.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Msg(err.Error())
	}

	log.Print("migrate up successfully")

}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start task processor")
	}
}
