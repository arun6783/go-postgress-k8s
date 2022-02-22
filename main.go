package main

import (
	"database/sql"
	"log"

	"github.com/arun6783/go-postgress-k8s/api"
	db "github.com/arun6783/go-postgress-k8s/db/sqlc"
	"github.com/arun6783/go-postgress-k8s/utils"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	var config utils.Config
	var conn *sql.DB

	config, err = utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Error occured when loading config")
	}
	conn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error occured when opening db connection", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Error occured when starting server", err)
	}
}
