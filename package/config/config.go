package config

import "fmt"

const DEFAULT_PORT = 6969

const DATABASE_PATH = "./database/"

var BASE_PATH = fmt.Sprintf("http://localhost:%d/v1", DEFAULT_PORT)
var IMAGE_PATH = fmt.Sprintf("%s/image", BASE_PATH)
var VIDEO_PATH = fmt.Sprintf("%s/video", BASE_PATH)
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
var VideoResolutions = map[string]string{
	"1080p": "5000k",
	"720p":  "2800k",
	"480p":  "1400k",
	"360p":  "800k",
	"240p":  "400k",
}

const UPLOAD_DIRECTORY = "./uploads"
const IMAGES_DIRECTORY = "images"
const VIDEOS_DIRECTORY = "videos"

var ImagesDirectory = fmt.Sprintf("%s/%s", UPLOAD_DIRECTORY, IMAGES_DIRECTORY)
var VideosDirectory = fmt.Sprintf("%s/%s", UPLOAD_DIRECTORY, IMAGES_DIRECTORY)

var ImageFormat = struct {
	WebP string
}{
	WebP: "webp",
}
