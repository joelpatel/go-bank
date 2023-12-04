package main

import (
	"fmt"

	"github.com/joelpatel/go-bank/db"
	"github.com/joelpatel/go-bank/db/controllers"
)

func main() {
	db.InitializeDBConnection()

	affected, err := controllers.DeleteAccountByID("1")
	if err != nil {
		panic(err)
	}

	fmt.Println(affected)
}
