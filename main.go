package main

import (
	"database/sql"
	"log"

	"github.com/arun6783/go-postgress-k8s/api"
	db "github.com/arun6783/go-postgress-k8s/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8086"
)

func main() {

	var conn *sql.DB
	var err error
	conn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Error occured when opening db connection", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Error occured when starting server", err)
	}
}
