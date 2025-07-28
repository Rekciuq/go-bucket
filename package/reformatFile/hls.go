package reformatfile

import (
	"log"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Rekciuq/go-bucket/package/config"
	"github.com/Rekciuq/go-bucket/package/ffmpeg"
	writefile "github.com/Rekciuq/go-bucket/package/writeFile"
)

func ConvertToHLS(videoFile multipart.File, outputDir string) (*ffmpeg.VideoMetadata, string, error) {
	err, inputFile := writefile.WriteTemporaryFile(videoFile, "upload-*.mp4")
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(inputFile.Name())

	metaData, err := ffmpeg.GetVideoMetadata(inputFile.Name())
	if err != nil {
		return nil, "", err
	}

	var targetResolutions []writefile.ResolutionInfo
	for resStr, bitrate := range config.VideoResolutions {
		height, _ := strconv.Atoi(strings.TrimSuffix(resStr, "p"))
		if height > 0 && height <= metaData.Height {
			targetResolutions = append(targetResolutions, writefile.ResolutionInfo{Height: height, Bitrate: bitrate})
		}
	}

	sort.Slice(targetResolutions, func(i, j int) bool {
		return targetResolutions[i].Height > targetResolutions[j].Height
	})
	log.Printf("Starting conversion for resolutions: %v", targetResolutions)

	var wg sync.WaitGroup
	errs := make(chan error, len(targetResolutions))

	for _, res := range targetResolutions {
		wg.Add(1)
		go func(r writefile.ResolutionInfo) {
			defer wg.Done()
			if err := ffmpeg.EncodeHLS(inputFile.Name(), outputDir, r); err != nil {
				errs <- err
			}
		}(res)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return nil, "", err
		}
	}

	masterPlaylist, err := writefile.WriteMasterPlaylist(outputDir, targetResolutions)
	if err != nil {
		return nil, "", err
	}
	log.Println("HLS conversion successful.")
	return metaData, masterPlaylist, nil
}
