package db

import (
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"

	"github.com/daniel/master-golang/utils"
)

var db *sql.DB

func InitDatabase(config utils.Config) (*sql.DB, error) {
	// mysql
	// dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?sslmode=disable",
	// 			config.Username,
	// 			config.Password,
	// 			config.Host,
	// 			config.Port,
	// 			config.Database,
	// )
	dsn := fmt.Sprintf("postgres://%s:%s@%v:%v/%s?sslmode=disable",
				config.Username,
				config.Password,
				config.Host,
				config.Port,
				config.Database,
	)
	fmt.Println(dsn)

	var err error
	db, err = sql.Open(config.Connection, dsn)

	return db, err
}

func Close() {
	db.Close()
}