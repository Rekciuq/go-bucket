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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
