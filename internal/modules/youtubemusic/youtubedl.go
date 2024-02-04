package youtubemusic

import (
	"log"
	"os/exec"
	"strings"
)

func GetYoutubeMusicUrl(videoID string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--get-url", "https://www.youtube.com/watch?v="+videoID)

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting youtube music url:", err)
		return "", err
	}

	url := strings.TrimSpace(string(output))
	return url, nil
}
