package writefile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Rekciuq/go-bucket/package/config"
)

func WriteImage(imageId string, webpBytes []byte) (string, error) {
	err := EnsureDir(config.ImagesDirectory)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}

	imageName := fmt.Sprintf("%s.%s", imageId, config.ImageFormat.WebP)
	imagePath := filepath.Join(config.ImagesDirectory, imageName)

	err = os.WriteFile(imagePath, webpBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return imagePath, nil

}
