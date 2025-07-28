package reformatfile

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os/exec"
)

func ConvertToWebP(imageFile multipart.File) ([]byte, error) {
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-f", "webp",
		"-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	defer imageFile.Close()
	data, err := io.ReadAll(imageFile)
	if err != nil {
		return nil, err
	}

	_, err = stdin.Write(data)
	if err != nil {
		stdin.Close()
		return nil, err
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %v, details: %s", err, errBuf.String())
	}

	return outBuf.Bytes(), nil
}
