package shared

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func GetRequestBody[T interface{}](r *http.Request) (T, error) {
	defer r.Body.Close()

	var body T

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, errors.New(ErrInvalidJSON)
	}

	return body, nil
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
