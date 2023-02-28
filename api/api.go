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

func handleStatusCode(code int) bool {
	if code < 200 || code >= 400 {
		return false
	}
	return true
}

func (us *Uspacy) doRaw(url, method string, headers map[string]string, body io.Reader) ([]byte, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", us.bearerToken)

	for key, value := range headers {
		req.Header.Add(key, value)
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

func (us *Uspacy) doGetEmptyHeaders(url string) ([]byte, error) {
	return us.doRaw(url, http.MethodGet, emptyHeaders, nil)
}

func buildURL(host, version, route string) string {
	return fmt.Sprintf("%s/%s/%s", host, version, route)
}
