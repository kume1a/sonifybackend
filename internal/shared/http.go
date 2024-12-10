package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/kume1a/sonifybackend/internal/config"
)

type Validatable interface {
	Validate() error
}

type HttpGetParams struct {
	URL     string
	Headers map[string]string
	Query   url.Values
}

func GetURLParamString(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)

	value, ok := vars[key]
	if !ok {
		return "", fmt.Errorf("missing URL parameter: %s", key)
	}

	return value, nil
}

func GetURLParamUUID(r *http.Request, key string) (uuid.UUID, error) {
	value, err := GetURLParamString(r, key)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %s", value)
	}

	return id, nil
}

func HttpGetWithResponse[DTO interface{}](params HttpGetParams) (*DTO, error) {
	req, err := http.NewRequest("GET", params.URL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}

	if len(params.Query) != 0 {
		req.URL.RawQuery = params.Query.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dto DTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

type XWWWFormUrlencodedParams struct {
	URL     string
	Form    url.Values
	Headers map[string]string
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

type HandleUploadFileArgs struct {
	ResponseWriter   http.ResponseWriter
	Request          *http.Request
	FieldName        string
	Dir              string
	AllowedMimeTypes []string
	IsOptional       bool
}

func HandleUploadFile(args HandleUploadFileArgs) (string, *HttpError) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", InternalServerErrorDef()
	}

	if args.Request.Method != "POST" {
		return "", MethodNotAllowed(ErrMethodNotAllowed)
	}

	args.Request.Body = http.MaxBytesReader(
		args.ResponseWriter,
		args.Request.Body,
		env.MaxUploadSizeBytes,
	)
	if err := args.Request.ParseMultipartForm(env.MaxUploadSizeBytes); err != nil {
		return "", BadRequest(ErrExceededMaxUploadSize)
	}

	file, fileHeader, err := args.Request.FormFile(args.FieldName)

	if err != nil {
		if err == http.ErrMissingFile && args.IsOptional {
			return "", nil
		}

		log.Println("error parsing form file: ", err)
		return "", BadRequest("field " + args.FieldName + " is required")
	}

	defer file.Close()

	if err := validateMimeType(file, args.AllowedMimeTypes); err != nil {
		return "", err
	}

	extension := filepath.Ext(fileHeader.Filename)
	location, err := NewFileLocation(FileLocationArgs{
		Dir:       args.Dir,
		Extension: extension,
	})
	if err != nil {
		return "", InternalServerErrorDef()
	}

	dst, err := os.Create(location)
	if err != nil {
		return "", InternalServerErrorDef()
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", InternalServerErrorDef()
	}

	return location, nil
}

func GetRequestBody[T interface{}](r *http.Request) (T, error) {
	defer r.Body.Close()

	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, errors.New(ErrInvalidJSON)
	}

	return body, nil
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
		return InternalServerErrorDef()
	}

	filetype := http.DetectContentType(buff)
	if !Contains(allowedMimeTypes, filetype) {
		log.Println("Invalid mime type: ", filetype, ", allowed = ", allowedMimeTypes)
		return BadRequest(ErrInvalidMimeType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return InternalServerErrorDef()
	}

	return nil
}
