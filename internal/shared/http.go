package shared

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Validatable interface {
	Validate() error
}

type XWWWFormUrlencodedParams struct {
	URL     string
	Form    url.Values
	Headers map[string]string
}

type HandleUploadFileArgs struct {
	ResponseWriter   http.ResponseWriter
	Request          *http.Request
	FieldName        string
	Dir              string
	AllowedMimeTypes []string
	IsOptional       bool
}

func GetRequestBody[T interface{}](r *http.Request) (T, error) {
	defer r.Body.Close()

	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, errors.New(ErrInvalidJSON)
	}

	return body, nil
}

func XWWWFormUrlencoded(params XWWWFormUrlencodedParams) (
	httpResp *http.Response,
	respBody string,
	err error,
) {
	req, err := http.NewRequest("POST", params.URL, strings.NewReader(params.Form.Encode()))
	if err != nil {
		return nil, "", err
	}

	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request: ", err)
		return nil, "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body: ", err)
		return nil, "", err
	}

	return resp, string(body), nil
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

func ValidateRequestBody[T Validatable](r *http.Request) (T, error) {
	defer r.Body.Close()

	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, errors.New(ErrInvalidJSON)
	}

	if err := body.Validate(); err != nil {
		log.Println(err)
		return body, err
	}

	return body, nil
}

func ValidateRequestQuery[T Validatable](r *http.Request) (T, error) {
	var q T

	jsonbody, err := json.Marshal(r.URL.Query())
	if err != nil {
		return q, errors.New(ErrInvalidJSON)
	}

	if err := json.Unmarshal(jsonbody, &q); err != nil {
		return q, errors.New(ErrInvalidJSON)
	}

	if err := q.Validate(); err != nil {
		return q, err
	}

	return q, nil
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
