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
	unixExpTime int64
	isExpired   bool
}

const defaultClientTimeout = 10 * time.Second

// New creates an Uspacy object
func New(token string) *Uspacy {
	return &Uspacy{
		bearerToken: token,
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
		unixExpTime: setExp(token),
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
	var token string

	if us.isExpired == true {
		token = us.bearerToken
	} else {
		token = us.getToken()
	}
	req.Header.Add("Authorization", token)

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

		return responseBody, errors.New("status code != 200")
	}

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

func buildURL(host, version, route string) string {
	return fmt.Sprintf("%s/%s/%s", host, version, route)
}
