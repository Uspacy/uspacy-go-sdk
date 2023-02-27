package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Uspacy struct {
	bearerToken string
	client      *http.Client
}

const defaultClientTimeout = 10 * time.Second

// New creates a Uspacy object
func New(token string) *Uspacy {
	return &Uspacy{
		bearerToken: "Bearer " + token,
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
	}
}

func (us *Uspacy) generateRequest(url, method string, headers map[string]string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", us.bearerToken)

	for key, value := range headers {
		request.Header.Add(key, value)
	}
	return request, nil
}

func handleStatusCode(code int) bool {
	if code < 200 || code >= 400 {
		return false
	}
	return true
}

func (us *Uspacy) doRaw(url, method string, headers map[string]string, body io.Reader) ([]byte, error) {

	req, err := us.generateRequest(mainHost+fmt.Sprintf(url), method, headers, body)
	if err != nil {
		return nil, err
	}

	res, err := us.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !handleStatusCode(res.StatusCode) {
		log.Printf("error occured while trying to (%s)\nbody - %s\ncode - %v\n", req.URL.String(), string(responseBody), res.StatusCode)
		return nil, err
	}

	return responseBody, nil

}
