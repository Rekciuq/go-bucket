package uploadhandler

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/Rekciuq/go-bucket/package/config"
	"github.com/go-playground/validator/v10"
)

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

type validatedForm struct {
	File    multipart.File
	Header  *multipart.FileHeader
	IsImage bool
}

func validateFormData(r *http.Request) (*validatedForm, error) {

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

	return &validatedForm{
		File:    file,
		Header:  handler,
		IsImage: isImage,
	}, nil
}
