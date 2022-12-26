package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB *sql.DB
)

const (
	dbDriver = "postgres"
	dbDns    = "postgres://daniel:mypassword@localhost:5432/test_db?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbDns)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
