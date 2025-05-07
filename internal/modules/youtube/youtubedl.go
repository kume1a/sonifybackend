package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type DownloadYoutubeAudioOptions struct {
	DownloadThumbnail bool
}

func DownloadYoutubeAudio(videoID string, options DownloadYoutubeAudioOptions) (outputPath string, thumbnailPath string, err error) {
	tempOutputDir := config.DirTempYoutubeAudios + "/" + videoID

	tempOutputPath, err := shared.NewFileLocation(shared.FileLocationArgs{
		Dir:       tempOutputDir,
		Extension: "",
	})
	if err != nil {
		return "", "", err
	}

	ytURL := "https://www.youtube.com/watch?v=" + videoID

	var cmd *exec.Cmd
	if options.DownloadThumbnail {
		cmd = exec.Command("yt-dlp", "-f", "bestaudio", "--write-thumbnail", "-o", tempOutputPath+".%(ext)s", ytURL)
	} else {
		cmd = exec.Command("yt-dlp", "-f", "bestaudio", "-o", tempOutputPath+".%(ext)s", ytURL)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		shared.LogCommandError(err, "DownloadYoutubeAudio")
		errorMessage := fmt.Sprintf("Error downloading youtube audio: %v, output: %s", err, output)
		return "", "", errors.New(errorMessage)
	}

	audioOutputLocation, err := shared.ProcessUnknownExtMediaFile(shared.ProcessUnknownExtMediaFileParams{
		TempOutputDir:      tempOutputDir,
		PossibleExtensions: shared.AudioExtensions,
		DesiredExtension:   shared.AudioExtMp3,
		OutputDir:          config.DirYoutubeAudios,
	})
	if err != nil {
		return "", "", err
	}

	thumbnailOutputLocation := ""
	if options.DownloadThumbnail {
		thumbnailOutputLocation, err = shared.ProcessUnknownExtMediaFile(shared.ProcessUnknownExtMediaFileParams{
			TempOutputDir:      tempOutputDir,
			PossibleExtensions: shared.ImageExtensions,
			DesiredExtension:   shared.ImageExtJpg,
			OutputDir:          config.DirYoutubeAudioThumbnails,
		})
		if err != nil {
			return "", "", err
		}
	}

	if err := os.RemoveAll(tempOutputDir); err != nil {
		log.Println("Error removing temp output directory: ", err)
		return "", "", err
	}

	return audioOutputLocation, thumbnailOutputLocation, nil
}

func GetYoutubeVideoInfo(videoID string) (*youtubeVideoInfoDTO, error) {
	formatString := `{"title": %(title)j, "uploader": %(uploader)j, "durationSeconds": %(duration)d}`

	cmd := exec.Command(
		"yt-dlp",
		"--print",
		formatString,
		"--skip-download",
		"--no-warnings",
		"https://www.youtube.com/watch?v="+videoID,
	)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("yt-dlp raw output on error: %s", string(output))
		shared.LogCommandError(err, "GetYoutubeVideoInfo")
		return nil, fmt.Errorf("yt-dlp execution failed: %w", err)
	}

	trimmedOutput := strings.TrimSpace(string(output))

	var info youtubeVideoInfoDTO
	if err := json.Unmarshal([]byte(trimmedOutput), &info); err != nil {
		log.Printf("Error parsing youtube video info from JSON: '%s'. Error: %v", trimmedOutput, err)
		return nil, fmt.Errorf("failed to parse yt-dlp JSON output: %w", err)
	}

	return &info, nil
}

func GetYoutubeAudioURL(query string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--get-url", "--default-search", "ytsearch", query)

	output, err := cmd.Output()
	if err != nil {
		shared.LogCommandError(err, "GetYoutubeAudioURL")
		return "", err
	}

	url := strings.TrimSpace(string(output))

	if url == "" {
		return "", errors.New("yt no audio found for query: " + query)
	}

	return url, nil
}

func GetYoutubeSearchBestMatchVideoID(query string) (string, error) {
	cmd := exec.Command("yt-dlp", "--get-id", "--default-search", "ytsearch", query)

	output, err := cmd.Output()
	if err != nil {
		shared.LogCommandError(err, "GetYoutubeSearchBestMatchVideoID")
		return "", err
	}

	ids := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(ids) == 0 {
		return "", errors.New("No video found for query: " + query)
	}

	return ids[0], nil
}
