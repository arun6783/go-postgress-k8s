package api

import (
	db "github.com/arun6783/go-postgress-k8s/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {

	server := &Server{store: store}

	router := gin.Default()

	server.bindRoutes(router)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) bindRoutes(router *gin.Engine) {

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

}
