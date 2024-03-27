package youtube

import (
	"encoding/json"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func GetYoutubeAudioUrl(videoID string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--get-url", "https://www.youtube.com/watch?v="+videoID)

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting youtube music url:", err)
		return "", err
	}

	url := strings.TrimSpace(string(output))
	return url, nil
}

func DownloadYoutubeAudio(videoID string) (outputPath string, thumbnailPath string, err error) {
	outputLocation, err := shared.NewPublicFileLocation(".webm")
	if err != nil {
		log.Println("Error creating public file location: ", err)
		return "", "", err
	}

	thumbnailLocation := strings.TrimSuffix(outputLocation, path.Ext(outputLocation)) + ".webp"

	var cmd *exec.Cmd
	if shared.IsPlatformLinux() {
		cmd = exec.Command("sudo", "yt-dlp", "-f", "bestaudio", "--write-thumbnail", "-o", outputLocation, "https://www.youtube.com/watch?v="+videoID)
	} else {
		cmd = exec.Command("yt-dlp", "-f", "bestaudio", "--write-thumbnail", "-o", outputLocation, "https://www.youtube.com/watch?v="+videoID)
	}

	if err := cmd.Run(); err != nil {
		log.Println("Error downloading youtube audio: ", err)
		return "", "", err
	}

	return outputLocation, thumbnailLocation, nil
}

func GetYoutubeVideoInfo(videoID string) (*YoutubeVideoInfo, error) {
	cmd := exec.Command(
		"yt-dlp",
		"--print",
		"{\"title\": \"%(title)s\", \"uploader\": \"%(uploader)s\", \"duration\": %(duration)s}",
		"https://www.youtube.com/watch?v="+videoID,
	)

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting youtube video info: ", err)
		return nil, err
	}

	var info YoutubeVideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		log.Println("Error parsing youtube video info: ", err)
		return nil, err
	}

	return &info, nil
}
