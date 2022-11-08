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

var testQueries *Queries

// Main entry point of all unit tests inside ONE specific golang PACKAGE.
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v\n", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
