package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/daniel/master-golang/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

func InitDatabase(config utils.Config) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s://%s:%s@%v:%v/%s?sslmode=disable",
				config.Connection,
				config.Username,
				config.Password,
				config.Host,
				config.Port,
				config.Database,
	)

	db, err = sql.Open(config.Connection, dsn)
	if err != nil {
		return
	}
	test(dsn)
	// err = initMigrate(
	// 	db,
	// 	config.Database,
	// 	config.Migration_url,
	// )

	return db, err
}

func Close() {
	db.Close()
}

func test(dbdsn string) {
	m, err := migrate.New(
        "file://db/migrations",
        dbdsn)

    if err != nil {
		log.Fatal(err.Error())
	}

	if err := m.Up(); err != nil  && err != migrate.ErrNoChange {
		log.Fatal(err.Error())
	}

	log.Println("success to migrate")
}

func initMigrate(db *sql.DB, database string, migrationUrl string) (err error) {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			migrationUrl,
			database, driver)
		m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

		if err != nil && err != migrate.ErrNoChange {
			return
		}

		fmt.Println("---------  initializing migrations is successful ----")
		return nil
}