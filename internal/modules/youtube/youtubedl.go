package youtube

import (
	"log"
	"os/exec"
	"strconv"
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

func DownloadYoutubeAudio(videoID string) (string, error) {
	outputLocation, err := shared.NewPublicFileLocation(".webm")
	if err != nil {
		log.Println("Error creating public file location: ", err)
		return "", err
	}

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", outputLocation, "https://www.youtube.com/watch?v="+videoID)

	if err := cmd.Run(); err != nil {
		log.Println("Error downloading youtube audio: ", err)
		return "", err
	}

	return outputLocation, nil
}

func GetYoutubeAudioDurationInSeconds(videoID string) (int, error) {
	cmd := exec.Command("yt-dlp", "--print", "%(duration>%s)s", "https://www.youtube.com/watch?v="+videoID)

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting youtube music duration:", err)
		return 0, err
	}

	duration, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		log.Println("Error converting duration to int:", err)
		return 0, err
	}

	return duration, nil
}

func GetYoutubeVideoTitle(videoID string) (string, error) {
	cmd := exec.Command("yt-dlp", "--print", "%(title)s", "https://www.youtube.com/watch?v="+videoID)

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting youtube video title:", err)
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
