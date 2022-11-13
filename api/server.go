package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/joelpatel/go-bank/db/sqlc"
)

// serves all HTTP requests for the banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// creates a new Server instance and setups up routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	// add routes to the router
	router.POST("/accounts/add", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)
	router.DELETE("/accounts/delete", server.deleteAccount)
	router.PUT("/accounts/update", server.updateAccount)
	router.POST("/transfer", server.createTranfer)

	server.router = router
	return server
}

// start an HTTP server on input address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
