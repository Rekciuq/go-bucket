package writefile

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func WriteTemporaryFile(file multipart.File, pattern string) (error, *os.File) {
	inputFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err), nil
	}

	_, err = io.Copy(inputFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy video to temp file: %w", err), nil
	}
	inputFile.Close()

	return nil, inputFile
}
