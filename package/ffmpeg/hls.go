package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	writefile "github.com/Rekciuq/go-bucket/package/writeFile"
)

func EncodeHLS(inputPath, outputDir string, res writefile.ResolutionInfo) error {
	resDirName := fmt.Sprintf("%dp", res.Height)
	resOutputDir := filepath.Join(outputDir, resDirName)
	if err := writefile.EnsureDir(resOutputDir); err != nil {
		return err
	}

	args := []string{
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=-2:%d", res.Height),
		"-c:a", "aac", "-b:a", "128k",
		"-c:v", "libx264", "-preset", "ultrafast", "-b:v", res.Bitrate,
		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(resOutputDir, "segment%03d.ts"),
		filepath.Join(resOutputDir, "playlist.m3u8"),
	}

	cmd := exec.Command("ffmpeg", args...)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error for %dp: %v, details: %s", res.Height, err, errBuf.String())
	}
	return nil
}
