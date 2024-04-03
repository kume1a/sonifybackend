package shared

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func HandleUploadFile(w http.ResponseWriter, r *http.Request, fieldName string, dir string, allowedMimeTypes []string) (string, *HttpError) {
	env, err := ParseEnv()
	if err != nil {
		return "", HttpErrInternalServerErrorDef()
	}

	if r.Method != "POST" {
		return "", HttpErrMethodNotAllowed(ErrMethodNotAllowed)
	}

	r.Body = http.MaxBytesReader(w, r.Body, env.MaxUploadSizeBytes)
	if err := r.ParseMultipartForm(env.MaxUploadSizeBytes); err != nil {
		return "", HttpErrBadRequest(ErrExceededMaxUploadSize)
	}

	file, fileHeader, err := r.FormFile(fieldName)

	if err != nil {
		log.Println("error parsing form file: ", err)
		return "", HttpErrBadRequest("field " + fieldName + " is required")
	}

	defer file.Close()

	if err := validateMimeType(file, allowedMimeTypes); err != nil {
		return "", err
	}

	extension := filepath.Ext(fileHeader.Filename)
	location, err := NewPublicFileLocation(PublicFileLocationArgs{
		Dir:       dir,
		Extension: extension,
	})
	if err != nil {
		return "", HttpErrInternalServerErrorDef()
	}

	dst, err := os.Create(location)
	if err != nil {
		return "", HttpErrInternalServerErrorDef()
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", HttpErrInternalServerErrorDef()
	}

	return location, nil
}

func validateMimeType(file multipart.File, allowedMimeTypes []string) *HttpError {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return HttpErrInternalServerErrorDef()
	}

	filetype := http.DetectContentType(buff)
	if !Contains(allowedMimeTypes, filetype) {
		return HttpErrBadRequest(ErrInvalidMimeType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return HttpErrInternalServerErrorDef()
	}

	return nil
}
