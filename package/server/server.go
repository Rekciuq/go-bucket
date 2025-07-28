package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rekciuq/go-bucket/ent"
	"github.com/Rekciuq/go-bucket/package/config"
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
	ih "github.com/Rekciuq/go-bucket/package/imageHandler"
	uploadhandler "github.com/Rekciuq/go-bucket/package/uploadHandler"
	uh "github.com/Rekciuq/go-bucket/package/userHandler"
)

func StartServer(ctx context.Context, client *ent.Client) *http.Server {
	router := http.NewServeMux()
	router.Handle("/v1/image/", http.StripPrefix("/v1/image", ih.ImageHandler(&entclient.ClientConnection{Ctx: ctx, Client: client})))
	router.Handle("/v1/users/", http.StripPrefix("/v1/users", uh.UserHandler(ctx, client)))
	router.Handle("/v1/upload/", http.StripPrefix("/v1/upload", uploadhandler.UploadHandler(&entclient.ClientConnection{Ctx: ctx, Client: client})))
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.DEFAULT_PORT),
		Handler: corsMiddleware(router),
	}

	fmt.Printf("Server listening on port %d\n", config.DEFAULT_PORT)
	return &server
}
