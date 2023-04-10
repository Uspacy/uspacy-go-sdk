package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Uspacy struct {
	bearerToken string
	client      *http.Client
	mainHost    string
	isExpired   bool
}

const defaultClientTimeout = 10 * time.Second

// New creates an Uspacy object
func New(token, host string) *Uspacy {
	return &Uspacy{
		bearerToken: token,
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
		mainHost: host,
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

	req.Header.Add("Authorization", us.getToken())
	req.Header.Add("Content-Type", "application/json")

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
		if res.StatusCode == http.StatusUnauthorized {
			if !us.isExpired {
				us.isExpired = true
				us.refreshToken()
				return us.doRaw(url, method, headers, body)
			}
		}
		log.Printf("error occured while trying to (%s)\nbody - %s\ncode - %v\n", req.URL.String(), string(responseBody), res.StatusCode)

		return responseBody, errors.New("status code != 200")
	}

	us.isExpired = false

	return responseBody, nil

}

func (us *Uspacy) doGetEmptyHeaders(url string) ([]byte, error) {
	return us.doRaw(url, http.MethodGet, emptyHeaders, nil)
}

func (us *Uspacy) doPostEmptyHeaders(url string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return us.doRaw(url, http.MethodPost, emptyHeaders, &buf)
}

func (us *Uspacy) doPatchEmptyHeaders(url string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return us.doRaw(url, http.MethodPatch, emptyHeaders, &buf)
}

func (us *Uspacy) buildURL(version, route string) string {
	return fmt.Sprintf("%s/%s/%s", us.mainHost, version, route)
}
