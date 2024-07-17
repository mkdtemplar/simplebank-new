package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	server.router = router

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
