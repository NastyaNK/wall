package main

import (
	"log"
	"wall/internal/repository"
	"wall/internal/repository/postgres"
	"wall/internal/web"
)

func main() {
	psql := postgres.NewDatabase()
	dbConfig, _ := repository.LoadConfig("./configs/db.json")
	err := psql.Connect(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	server := web.NewServer(psql)
	log.Fatal(server.Run())
}
