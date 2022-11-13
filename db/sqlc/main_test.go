// Main test file.
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joelpatel/go-bank/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

// Main entry point of all unit tests inside ONE specific golang PACKAGE.
func TestMain(m *testing.M) {
	config := util.LoadConfig()
	var err error
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v\n", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
