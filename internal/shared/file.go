package shared

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileSize struct {
	Bytes     int64
	KiloBytes float64
	MegaBytes float64
}

type FileLocationArgs struct {
	Extension string
	Dir       string
}

func NewFileLocation(args FileLocationArgs) (string, error) {
	if err := ensureDir(args.Dir); err != nil {
		log.Println("error ensuring dir: ", err)
		return "", err
	}

	fileName := uuid.New().String()
	if args.Extension != "" {
		fileName = fileName + "." + strings.Trim(args.Extension, ".")
	}

	return args.Dir + "/" + fileName, nil
}

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		log.Println("error creating file: ", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error getting url: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("bad status: ", resp.Status)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		log.Println("error copying file: ", err)
		return err
	}

	return nil
}

func GetFileSize(filepath string) (*FileSize, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		log.Printf("error getting file info: %v", err)
		return nil, err
	}

	return &FileSize{
		Bytes:     fileInfo.Size(),
		KiloBytes: float64(fileInfo.Size()) / 1024,
		MegaBytes: float64(fileInfo.Size()) / (1024 * 1024),
	}, nil
}

func DeleteFiles(filepaths []string) error {
	for _, filepath := range filepaths {
		if err := os.Remove(filepath); err != nil {
			log.Println("error deleting file: ", err)
			return err
		}
	}
	return nil
}

func ReplaceFilenameExtension(filename string, newExtension string) string {
	extension := filepath.Ext(filename)

	return strings.TrimSuffix(filename, extension) + newExtension
}

func MoveFile(oldPath string, newPath string) error {
	err := os.Rename(oldPath, newPath)

	if err != nil {
		log.Println(err)
	}

	return err
}

func FindFirstFile(dir string, extensions []string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		for _, ext := range extensions {
			fileName := file.Name()

			if strings.HasSuffix(fileName, ext) {
				return fileName, nil
			}
		}
	}

	return "", fmt.Errorf("no file with ext: %s, found in directory: %s", extensions, dir)
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err == nil {
		return nil
	}

	if os.IsExist(err) {
		return nil
	}

	info, err := os.Stat(dirName)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("path exists but is not a directory")
	}
	return nil
}
