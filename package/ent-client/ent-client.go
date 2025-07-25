package entclient

import (
	"context"
	"fmt"
	"log"

	"github.com/Rekciuq/go-bucket/ent"
	_ "github.com/mattn/go-sqlite3"
)

const databasePath = "./database/db.db?_fk=1"

type ClientConnection struct {
	Ctx    context.Context
	Client *ent.Client
}

func CreateClient(ctx context.Context) *ent.Client {
	client, err := ent.Open("sqlite3", fmt.Sprintf("file:%s", databasePath))
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}
