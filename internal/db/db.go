package db

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"github.com/teodora-lazarevic/Poll-App/ent"
)

func InitDB(dbUrl string) *ent.Client {
	client, err := ent.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalf("Failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed creating schema resources: %v", err)
	}

	log.Println("Database connection and schema migrations successful")
	return client
}
