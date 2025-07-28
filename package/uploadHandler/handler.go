package uploadhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Rekciuq/go-bucket/package/config"
	ec "github.com/Rekciuq/go-bucket/package/ent-client"
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
	reformatfile "github.com/Rekciuq/go-bucket/package/reformatFile"
	writefile "github.com/Rekciuq/go-bucket/package/writeFile"
	"github.com/google/uuid"
)

type handler struct {
	*entclient.ClientConnection
	*uploadController
}

func generateUUID() uuid.UUID {
	uuID, err := uuid.NewRandom()
	if err != nil {
		generateUUID()
	}

	return uuID
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	uuID := generateUUID()
	isExists := h.uploadController.IsAlreadyExists(uuID.String())
	if isExists {
		uuID = generateUUID()
	}

	err := h.uploadController.Create(uuID.String())

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GET: %s", "success")
	res := getPayload{URL: fmt.Sprintf("%s/%s", config.UPLOAD_PATH, uuID.String())}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(res)
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	urlId := strings.TrimPrefix(r.URL.Path, "/")

	url, err := h.uploadController.GetUrl(urlId)
	isExpired := url.ExpiresAt.Before(time.Now())

	if err != nil || url.IsUsed {
		log.Println(err)
		http.Error(w, "Url is used!", http.StatusBadRequest)
		return
	}

	if isExpired {
		expiredUrlError := "Url is expired!"
		log.Println(expiredUrlError)
		http.Error(w, expiredUrlError, http.StatusBadRequest)
		return
	}

	validatedData, err := validateFormData(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer validatedData.File.Close()

	log.Printf("Received valid data: UrlID=%s, Filename=%s, IsImage=%t", urlId, validatedData.Header.Filename, validatedData.IsImage)

	if validatedData.IsImage {
		webpBytes, err := reformatfile.ConvertToWebP(validatedData.File)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		imagePath, err := writefile.WriteImage(urlId, webpBytes)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.uploadController.CreateImage(urlId, imagePath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.uploadController.UseUrl(urlId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		res := postResponse{ImageURL: fmt.Sprintf("%s/%s", config.IMAGE_PATH, urlId)}
		json.NewEncoder(w).Encode(res)
	}

	if !validatedData.IsImage {
		resolutions, err := reformatfile.ConvertToHLS(validatedData.File, fmt.Sprintf("./uploads/videos/%s", urlId))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println(resolutions)

	}
}

func UploadHandler(connection *ec.ClientConnection) *http.ServeMux {
	router := http.NewServeMux()
	h := handler{connection, &uploadController{connection: connection}}

	router.HandleFunc("GET /", h.get)
	router.HandleFunc("POST /{id}", h.post)
	router.HandleFunc("PUT /{id}", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {})

	return router
}
