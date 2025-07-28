package writefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ResolutionInfo struct {
	Height  int
	Bitrate string
}

func WriteMasterPlaylist(outputDir string, resolutions []ResolutionInfo) (string, error) {
	masterPlaylistPath := filepath.Join(outputDir, "master.m3u8")
	content := "#EXTM3U\n#EXT-X-VERSION:3\n"

	for _, res := range resolutions {
		resDirName := fmt.Sprintf("%dp", res.Height)
		bandwidth := 6000000
		if val, err := strconv.Atoi(strings.TrimSuffix(res.Bitrate, "k")); err == nil {
			bandwidth = val * 1000
		}

		width := (res.Height * 16) / 9
		content += fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n", bandwidth, width, res.Height)
		content += fmt.Sprintf("%s/playlist.m3u8\n", resDirName)
	}

	return masterPlaylistPath, os.WriteFile(masterPlaylistPath, []byte(content), 0644)
}
