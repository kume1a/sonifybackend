package youtube

import (
	"encoding/json"
	"errors"
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

func DownloadYoutubeAudioWithThumbnail(videoID string) (outputPath string, thumbnailPath string, err error) {
	outputLocation, err := shared.NewPublicFileLocation(shared.PublicFileLocationArgs{
		Dir:       shared.DirPublic,
		Extension: "mp3",
	})
	if err != nil {
		log.Println("Error creating public file location: ", err)
		return "", "", err
	}

	thumbnailLocation := strings.TrimSuffix(outputLocation, path.Ext(outputLocation)) + ".webp"

	ytURL := "https://www.youtube.com/watch?v=" + videoID

	var cmd *exec.Cmd
	if shared.IsPlatformLinux() {
		cmd = exec.Command("sudo", "yt-dlp", "-f", "bestaudio", "--write-thumbnail", "-o", outputLocation, ytURL)
	} else {
		cmd = exec.Command("yt-dlp", "-f", "bestaudio", "--write-thumbnail", "-o", outputLocation, ytURL)
	}

	if err := cmd.Run(); err != nil {
		log.Println("Error downloading youtube audio: ", err)
		return "", "", err
	}

	return outputLocation, thumbnailLocation, nil
}

func DownloadYoutubeAudio(videoID, outputPath string) (err error) {
	ytURL := "https://www.youtube.com/watch?v=" + videoID

	var cmd *exec.Cmd
	if shared.IsPlatformLinux() {
		cmd = exec.Command("sudo", "yt-dlp", "-f", "bestaudio", "-o", outputPath, ytURL)
	} else {
		cmd = exec.Command("yt-dlp", "-f", "bestaudio", "-o", outputPath, ytURL)
	}

	if err := cmd.Run(); err != nil {
		log.Println("Error downloading youtube audio: ", err)
		return err
	}

	return nil
}

func GetYoutubeVideoInfo(videoID string) (*youtubeVideoInfoDTO, error) {
	cmd := exec.Command(
		"yt-dlp",
		"--print",
		"{\"title\": \"%(title)s\", \"uploader\": \"%(uploader)s\", \"durationSeconds\": %(duration)s}",
		"https://www.youtube.com/watch?v="+videoID,
	)

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting youtube video info: ", err)
		return nil, err
	}

	var info youtubeVideoInfoDTO
	if err := json.Unmarshal(output, &info); err != nil {
		log.Println("Error parsing youtube video info: ", err)
		return nil, err
	}

	return &info, nil
}

func GetYoutubeAudioURL(query string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--get-url", "--default-search", "ytsearch", "\""+query+"\"")

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error searching youtube music video:", err)
		return "", err
	}

	url := strings.TrimSpace(string(output))

	if url == "" {
		return "", errors.New("yt no audio found for query: " + query)
	}

	return url, nil
}

func GetYoutubeSearchBestMatchVideoID(query string) (string, error) {
	cmd := exec.Command("yt-dlp", "--get-id", "--default-search", "ytsearch", "\""+query+"\"")

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error searching youtube: ", err)
		return "", err
	}

	ids := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(ids) == 0 {
		return "", errors.New("No video found for query: " + query)
	}

	return ids[0], nil
}
