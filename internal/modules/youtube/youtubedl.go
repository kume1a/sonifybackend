package youtube

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
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

type DownloadYoutubeAudioOptions struct {
	DownloadThumbnail bool
}

func DownloadYoutubeAudio(videoID string, options DownloadYoutubeAudioOptions) (outputPath string, thumbnailPath string, err error) {
	tempOutputDirectory := shared.DirTempYoutubeAudios + "/" + videoID

	tempOutputPath, err := shared.NewFileLocation(shared.FileLocationArgs{
		Dir:       tempOutputDirectory,
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

	if err := cmd.Run(); err != nil {
		log.Println("Error downloading youtube audio: ", err)
		return "", "", err
	}

	downloadedAudioFileName, err := shared.FindFirstFile(tempOutputDirectory, shared.AudioExtensions)
	if err != nil {
		return "", "", err
	}

	audioAsMp3Path := tempOutputDirectory + "/" + shared.ReplaceFilenameExtension(downloadedAudioFileName, shared.AudioExtMp3)
	nonMp3AudioPath := tempOutputDirectory + "/" + downloadedAudioFileName
	if !strings.HasSuffix(downloadedAudioFileName, shared.AudioExtMp3) {
		if err := shared.ConvertMedia(nonMp3AudioPath, audioAsMp3Path); err != nil {
			return "", "", err
		}

		if err := os.Remove(nonMp3AudioPath); err != nil {
			log.Println("Error removing non-mp3 audio: ", err)
			return "", "", err
		}
	}

	audioOutputLocation, err := shared.NewFileLocation(shared.FileLocationArgs{
		Dir:       shared.DirYoutubeAudios,
		Extension: shared.AudioExtMp3,
	})
	if err != nil {
		return "", "", err
	}

	if err := shared.MoveFile(audioAsMp3Path, audioOutputLocation); err != nil {
		return "", "", err
	}

	thumbnailOutputLocation := ""
	if options.DownloadThumbnail {
		downloadedThumbnailFileName, err := shared.FindFirstFile(tempOutputDirectory, shared.ImageExtensions)
		if err != nil {
			return "", "", err
		}

		thumbnailAsJpgPath := tempOutputDirectory + "/" + shared.ReplaceFilenameExtension(downloadedThumbnailFileName, shared.ImageExtJpg)
		nonJpgThumbnailPath := tempOutputDirectory + "/" + downloadedThumbnailFileName
		if !strings.HasSuffix(downloadedThumbnailFileName, shared.ImageExtJpg) {
			if err := shared.ConvertMedia(nonJpgThumbnailPath, thumbnailAsJpgPath); err != nil {
				return "", "", err
			}

			if err := os.Remove(nonJpgThumbnailPath); err != nil {
				log.Println("Error removing non-jpg thumbnail: ", err)
				return "", "", err
			}
		}

		thumbnailOutputLocation, err = shared.NewFileLocation(shared.FileLocationArgs{
			Dir:       shared.DirYoutubeAudioThumbnails,
			Extension: shared.ImageExtJpg,
		})
		if err != nil {
			return "", "", err
		}

		if err := shared.MoveFile(thumbnailAsJpgPath, thumbnailOutputLocation); err != nil {
			return "", "", err
		}
	}

	if err := os.Remove(tempOutputDirectory); err != nil {
		log.Println("Error removing temp output directory: ", err)
		return "", "", err
	}

	return audioOutputLocation, thumbnailOutputLocation, nil
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
