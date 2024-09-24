package gapi

import (
	"fmt"

	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
	"github.com/mkdtemplar/simplebank-new/pb"
	"github.com/mkdtemplar/simplebank-new/token"
	"github.com/mkdtemplar/simplebank-new/util"
	"github.com/mkdtemplar/simplebank-new/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %v", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config, taskDistributor: taskDistributor}

	return server, nil
}
