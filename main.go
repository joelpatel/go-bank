package main

import (
	"log"

	"github.com/joelpatel/go-bank/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db.InitializeDBConnection()
}
