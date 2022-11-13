package main

import (
	"database/sql"
	"log"

	"github.com/joelpatel/go-bank/api"
	db "github.com/joelpatel/go-bank/db/sqlc"
	"github.com/joelpatel/go-bank/util"
	_ "github.com/lib/pq"
)

// Main entry point of all unit tests inside ONE specific golang PACKAGE.
func main() {
	config := util.LoadConfig()
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v\n", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
