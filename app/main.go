package main

import (
	"log"

	"github.com/daniel/master-golang/api"
	db "github.com/daniel/master-golang/db"
	sqlc "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/utils"
)

var ENV_PATH = ".env"

func main() {
	config, err := utils.LoadConfig(ENV_PATH)
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.InitDatabase(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := sqlc.NewStore(db)
	router := api.NewServer(store)
	log.Fatal(router.Run())
}