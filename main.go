package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mkdtemplar/simplebank-new/api"
	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
	"github.com/mkdtemplar/simplebank-new/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load congiguration %s", err.Error())
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
