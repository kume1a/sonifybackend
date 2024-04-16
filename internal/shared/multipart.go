package shared

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type HandleUploadFileArgs struct {
	ResponseWriter   http.ResponseWriter
	Request          *http.Request
	FieldName        string
	Dir              string
	AllowedMimeTypes []string
	IsOptional       bool
}

func HandleUploadFile(args HandleUploadFileArgs) (string, *HttpError) {
	env, err := ParseEnv()
	if err != nil {
		return "", HttpErrInternalServerErrorDef()
	}

	if args.Request.Method != "POST" {
		return "", HttpErrMethodNotAllowed(ErrMethodNotAllowed)
	}

	args.Request.Body = http.MaxBytesReader(args.ResponseWriter, args.Request.Body, env.MaxUploadSizeBytes)
	if err := args.Request.ParseMultipartForm(env.MaxUploadSizeBytes); err != nil {
		return "", HttpErrBadRequest(ErrExceededMaxUploadSize)
	}

	file, fileHeader, err := args.Request.FormFile(args.FieldName)

	if err != nil {
		if err == http.ErrMissingFile && args.IsOptional {
			return "", nil
		}

		log.Println("error parsing form file: ", err)
		return "", HttpErrBadRequest("field " + args.FieldName + " is required")
	}

	defer file.Close()

	if err := validateMimeType(file, args.AllowedMimeTypes); err != nil {
		return "", err
	}

	extension := filepath.Ext(fileHeader.Filename)
	location, err := NewPublicFileLocation(PublicFileLocationArgs{
		Dir:       args.Dir,
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
		log.Println("Invalid mime type: ", filetype, ", allowed = ", allowedMimeTypes)
		return HttpErrBadRequest(ErrInvalidMimeType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return HttpErrInternalServerErrorDef()
	}

	return nil
}
