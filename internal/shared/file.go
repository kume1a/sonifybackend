package shared

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type FileSize struct {
	Bytes     int64
	KiloBytes float64
	MegaBytes float64
}

type PublicFileLocationArgs struct {
	Extension string
	Dir       string
}

func NewPublicFileLocation(args PublicFileLocationArgs) (string, error) {
	if err := ensureDir(args.Dir); err != nil {
		log.Println("error ensuring dir: ", err)
		return "", err
	}

	fileName := uuid.New().String()
	if args.Extension != "" {
		fileName = fileName + "." + args.Extension
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
