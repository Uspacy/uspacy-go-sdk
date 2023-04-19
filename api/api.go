package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

	if len(headers) == 0 {
		req.Header.Add("Content-Type", "application/json")
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
		if res.StatusCode == http.StatusUnauthorized {
			if !us.isExpired {
				us.isExpired = true
				if us.refreshToken() == nil {
					return us.doRaw(url, method, headers, body)
				} else {
					return nil, err
				}
			}
		}
		errMsg := fmt.Sprintf("error occured while trying to (%s)\nbody - %s\ncode - %v\n", req.URL.String(), string(responseBody), res.StatusCode)

		log.Println(errMsg)
		return responseBody, errors.New(string(responseBody))
	}
	us.isExpired = false
	return responseBody, nil
}

func (us *Uspacy) doGetEmptyHeaders(url string) ([]byte, error) {
	return us.doRaw(url, http.MethodGet, headersMap, nil)
}

func (us *Uspacy) doPostEmptyHeaders(url string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return us.doRaw(url, http.MethodPost, headersMap, &buf)
}

func (us *Uspacy) doPatchEmptyHeaders(url string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return us.doRaw(url, http.MethodPatch, headersMap, &buf)
}
func (us *Uspacy) doPostFormData(url string, body url.Values) ([]byte, error) {
	var buf bytes.Buffer
	headersMap["Content-Type"] = "application/x-www-form-urlencoded"
	headersMap["Accept"] = "application/json"
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return us.doRaw(url, http.MethodPost, headersMap, &buf)
}

func (us *Uspacy) buildURL(version, route string) string {
	return fmt.Sprintf("%s/%s/%s", us.mainHost, version, route)
}
