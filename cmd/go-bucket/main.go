package main

import (
	"context"
	"log"

	"github.com/Rekciuq/go-bucket/package/config"
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
	"github.com/Rekciuq/go-bucket/package/server"
	writefile "github.com/Rekciuq/go-bucket/package/writeFile"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	writefile.EnsureDir(config.DATABASE_PATH)
	ctx := context.Background()
	client := entclient.CreateClient(ctx)
	defer client.Close()

	server := server.StartServer(ctx, client)
	defer server.Close()

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
