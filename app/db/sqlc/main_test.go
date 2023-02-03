package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/daniel/master-golang/utils"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB *sql.DB
)


func TestMain(m *testing.M) {
	var err error

	config, err := utils.LoadConfig("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s://%s:%s@%v:%v/%s?sslmode=disable",
		config.Connection,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	testDB, err = sql.Open(config.Connection, dsn)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
