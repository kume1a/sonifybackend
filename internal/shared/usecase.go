package shared

import (
	"log"
	"os"
	"strings"
)

type ProcessUnknownExtMediaFileParams struct {
	TempOutputDir      string
	PossibleExtensions []string
	DesiredExtension   string
	OutputDir          string
}

func ProcessUnknownExtMediaFile(params ProcessUnknownExtMediaFileParams) (string, error) {
	downloadedFileName, err := FindFirstFile(params.TempOutputDir, params.PossibleExtensions)
	if err != nil {
		return "", err
	}

	fileAsDesiredExtPath := params.TempOutputDir + "/" + ReplaceFilenameExtension(downloadedFileName, params.DesiredExtension)
	nonDesiredExtFilePath := params.TempOutputDir + "/" + downloadedFileName
	if !strings.HasSuffix(downloadedFileName, params.DesiredExtension) {
		if err := ConvertMedia(nonDesiredExtFilePath, fileAsDesiredExtPath); err != nil {
			return "", err
		}

		if err := os.Remove(nonDesiredExtFilePath); err != nil {
			log.Println("Error removing non-desired ext file: ", err)
			return "", err
		}
	}

	outputLocation, err := NewFileLocation(FileLocationArgs{
		Dir:       params.OutputDir,
		Extension: params.DesiredExtension,
	})
	if err != nil {
		return "", err
	}

	if err := MoveFile(fileAsDesiredExtPath, outputLocation); err != nil {
		return "", err
	}

	return outputLocation, nil
}
