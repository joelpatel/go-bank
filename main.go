package main

import (
	"log"

	"github.com/joelpatel/go-bank/db"
	"github.com/joho/godotenv"
)

var store db.Store

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	store = db.InitializeDBStore()
}
