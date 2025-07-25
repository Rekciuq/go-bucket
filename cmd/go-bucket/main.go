package main

import (
	"context"

	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
	"github.com/Rekciuq/go-bucket/package/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()
	client := entclient.CreateClient(ctx)
	defer client.Close()

	server := server.StartServer(ctx, client)
	defer server.Close()

	server.ListenAndServe()
}
