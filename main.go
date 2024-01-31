package main

import (
	"log"
	"os"

	"github.com/joelpatel/go-bank/api"
	"github.com/joelpatel/go-bank/db"
	"github.com/joho/godotenv"
)

var (
	store  db.Store
	server *api.Server
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")

	store = db.InitializeDBStore()
	server = api.NewServer(store)

	err = server.StartServer(serverAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}
