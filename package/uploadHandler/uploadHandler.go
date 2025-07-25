package uploadhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/Rekciuq/go-bucket/package/config"
	ec "github.com/Rekciuq/go-bucket/package/ent-client"
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
	reformatfile "github.com/Rekciuq/go-bucket/package/reformatFile"
	uploadcontroller "github.com/Rekciuq/go-bucket/package/uploadController"
	writefile "github.com/Rekciuq/go-bucket/package/writeFile"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type handler struct {
	*entclient.ClientConnection
	*uploadcontroller.UploadController
}

func generateUUID() uuid.UUID {
	uuID, err := uuid.NewRandom()
	if err != nil {
		generateUUID()
	}

	return uuID
}

type getPayload struct {
	URL string `json:"url"`
}

type postPayload struct {
	UserID int `validate:"required,gt=0"`
}

type postResponse struct {
	ImageURL string `json:"imageURL"`
}

var validate = validator.New()

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	uuID := generateUUID()
	isExists := h.UploadController.IsAlreadyExists(uuID.String())
	if isExists {
		uuID = generateUUID()
	}

	err := h.UploadController.Create(uuID.String())

	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create url, please try again", http.StatusInternalServerError)
		return
	}

	log.Printf("GET: %s", "success")
	res := getPayload{URL: fmt.Sprintf("%s/%s", config.UPLOAD_PATH, uuID.String())}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(res)
}

type ValidatedForm struct {
	File    multipart.File
	Header  *multipart.FileHeader
	IsImage bool
}

func validateFormData(r *http.Request) (*ValidatedForm, error) {

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return nil, errors.New("File is missing or invalid")
	}

	contentType := handler.Header.Get("Content-Type")

	var isImage bool
	if _, ok := config.ImageTypes[contentType]; ok {
		isImage = true
	} else if _, ok := config.VideoTypes[contentType]; ok {
		isImage = false
	} else {
		file.Close()
		return nil, fmt.Errorf("invalid file type: %s", contentType)
	}

	return &ValidatedForm{
		File:    file,
		Header:  handler,
		IsImage: isImage,
	}, nil
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	urlId := strings.TrimPrefix(r.URL.Path, "/")

	url, err := h.UploadController.GetUrl(urlId)
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

		err = h.UploadController.CreateImage(urlId, imagePath)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.UploadController.UseUrl(urlId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		res := postResponse{ImageURL: fmt.Sprintf("%s/%s", config.IMAGE_PATH, urlId)}
		json.NewEncoder(w).Encode(res)
	}
}

func UploadHandler(connection *ec.ClientConnection) *http.ServeMux {
	router := http.NewServeMux()
	h := handler{connection, &uploadcontroller.UploadController{Connection: connection}}

	router.HandleFunc("GET /", h.get)
	router.HandleFunc("POST /{id}", h.post)
	router.HandleFunc("PUT /{id}", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {})

	return router
}
