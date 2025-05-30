package shared

import (
	"os/exec"
)

const (
	ImageExtJpg  = ".jpg"
	ImageExtJpeg = ".jpeg"
	ImageExtPng  = ".png"
	ImageExtBmp  = ".bmp"
	ImageExtTiff = ".tiff"
	ImageExtWebp = ".webp"
)

const (
	AudioExtMp3  = ".mp3"
	AudioExtWav  = ".wav"
	AudioExtAiff = ".aiff"
	AudioExtFlac = ".flac"
	AudioExtOgg  = ".ogg"
	AudioExtWma  = ".wma"
	AudioExtWmv  = ".wmv"
	AudioExtWebm = ".webm"
	AudioExtMp4  = ".mp4"
)

var AudioMimeTypes = []string{
	"audio/mpeg", "audio/mp3", "audio/x-m4a",
	"audio/x-wav", "audio/x-aiff", "audio/x-flac",
	"audio/ogg", "audio/x-ms-wma", "audio/x-ms-wmv",
	"audio/x-ms-wav", "audio/webm", "video/webm",
}

var ImageMimeTypes = []string{"image/jpeg", "image/png"}

var ImageExtensions = []string{
	ImageExtJpeg, ImageExtJpg, ImageExtPng,
	ImageExtBmp, ImageExtTiff, ImageExtWebp,
}

var AudioExtensions = []string{
	AudioExtMp3, AudioExtWav, AudioExtAiff,
	AudioExtFlac, AudioExtOgg, AudioExtWma,
	AudioExtWmv, AudioExtWebm, AudioExtMp4,
}

func ConvertMedia(inputFile string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, outputFile)

	if err := cmd.Run(); err != nil {
		LogCommandError(err, "ConvertMedia")
		return err
	}

	return nil
}
