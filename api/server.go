package api

import (
	"github.com/gin-gonic/gin"
	"github.com/joelpatel/go-bank/db"
)

// Server serves HTTP requests for the banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server instance and sets up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/account/create", server.createAccount)
	router.GET("/account/:id", server.getAccountByID)
	router.POST("/accounts", server.listAccountsByOwner)
	router.PUT("/account/update", server.updateAccountOwner)
	router.DELETE("/account/delete/:id", server.deleteAccountByID)

	server.router = router
	return server
}

// StartServer runs the HTTP server on a provided address.
func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
