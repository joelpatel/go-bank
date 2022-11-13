// Main test file.
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

// Main entry point of all unit tests inside ONE specific golang PACKAGE.
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v\n", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
