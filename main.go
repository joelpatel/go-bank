package main

import (
	"database/sql"
	"log"

	"github.com/joelpatel/go-bank/api"
	db "github.com/joelpatel/go-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:password@localhost:5432/bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

// Main entry point of all unit tests inside ONE specific golang PACKAGE.
func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v\n", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
