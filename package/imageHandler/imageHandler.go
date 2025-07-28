package imagehandler

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ec "github.com/Rekciuq/go-bucket/package/ent-client"
)

type handler struct {
	*ImageController
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	imageId := strings.TrimPrefix(r.URL.Path, "/")
	imagePath, err := h.GetImage(imageId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("image path", imagePath)

	image, err := os.Open(imagePath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer image.Close()

	ext := filepath.Ext(imagePath)
	mimeType := mime.TypeByExtension(ext)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, image)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ImageHandler(connection *ec.ClientConnection) http.Handler {
	router := http.NewServeMux()

	h := handler{&ImageController{Connection: connection}}
	router.HandleFunc("GET /{id}", h.get)

	router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("It's post!")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read the body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		fmt.Println(string(body))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Response!")
	})

	return router
}
