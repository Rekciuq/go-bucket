package config

import "fmt"

const DEFAULT_PORT = 6969

const DATABASE_PATH = "./database/"

var BASE_PATH = fmt.Sprintf("http://localhost:%d/v1", DEFAULT_PORT)
var IMAGE_PATH = fmt.Sprintf("%s/image", BASE_PATH)
var UPLOAD_PATH = fmt.Sprintf("%s/upload", BASE_PATH)

var ImageTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

var VideoTypes = map[string]struct{}{
	"video/mp4":       {},
	"video/webm":      {},
	"video/quicktime": {},
}

const UPLOAD_DIRECTORY = "./uploads"
const IMAGES_DIRECTORY = "images"

var ImagesDirectory = fmt.Sprintf("%s/%s", UPLOAD_DIRECTORY, IMAGES_DIRECTORY)

var ImageFormat = struct {
	WebP string
}{
	WebP: "webp",
}
