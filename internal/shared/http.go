package shared

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type XWWWFormUrlencodedParams struct {
	URL     string
	Form    url.Values
	Headers map[string]string
}

func XWWWFormUrlencoded[RESPONSE interface{}](params XWWWFormUrlencodedParams) (*RESPONSE, error) {
	req, err := http.NewRequest("POST", params.URL, strings.NewReader(params.Form.Encode()))
	if err != nil {
		return nil, err
	}

	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body: ", err)
		return nil, err
	}

	var dto RESPONSE
	if err := json.Unmarshal(body, &dto); err != nil {
		log.Println("error unmarshalling response body: ", err)
		return nil, err
	}

	return &dto, nil
}
