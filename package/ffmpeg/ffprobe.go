package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type stream struct {
	CodecType          string `json:"codec_type"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
}

type format struct {
	Duration string `json:"duration"`
}

type fFProbeOutput struct {
	Streams []stream `json:"streams"`
	Format  format   `json:"format"`
}

type VideoMetadata struct {
	Width           int
	Height          int
	DurationSeconds float64
	AspectRatio     string
}

func GetVideoMetadata(filePath string) (*VideoMetadata, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	)

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w, details: %s", err, errBuf.String())
	}

	var probeData fFProbeOutput
	if err := json.Unmarshal(output, &probeData); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	duration, err := strconv.ParseFloat(probeData.Format.Duration, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	for _, stream := range probeData.Streams {
		if stream.CodecType == "video" {
			return &VideoMetadata{
				Width:           stream.Width,
				Height:          stream.Height,
				DurationSeconds: duration,
				AspectRatio:     stream.DisplayAspectRatio,
			}, nil
		}
	}

	return nil, fmt.Errorf("no video stream found in file")
}
