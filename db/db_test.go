package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var testConn Executor

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	host := os.Getenv("TEST_DATABASE_HOST")
	port := os.Getenv("TEST_DATABASE_PORT")
	user := os.Getenv("TEST_DATABASE_USER")
	pass := os.Getenv("TEST_DATABASE_PASS")
	dbName := os.Getenv("TEST_DATABASE_NAME")
	sslMode := os.Getenv("TEST_DATABASE_SSLMODE")

	databaseUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbName, sslMode)

	db, err := sqlx.Open("pgx", databaseUrl)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err.Error())
	}

	testConn = db

	os.Exit(m.Run())
}
