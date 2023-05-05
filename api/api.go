package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	//"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Uspacy struct {
	bearerToken  string
	refreshToken string
	client       *http.Client
	mainHost     string
	isExpired    bool
	retries      int `default0:"3"`
}

const (
	defaultClientTimeout = 10 * time.Second
	defaultRetries       = 3
	tokenPrefix          = "Bearer "
)

// New creates an Uspacy object
func New(token, refresh, host string) *Uspacy {

	return &Uspacy{
		bearerToken:  strings.TrimPrefix(token, tokenPrefix),
		refreshToken: strings.TrimPrefix(refresh, tokenPrefix),
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
		mainHost:  host,
		retries:   defaultRetries,
		isExpired: false,
	}
}

// if WithRetries not set, default value will be 3
func (us *Uspacy) WithRetries(retries int) *Uspacy {
	us.retries = retries
	return us
}

func handleStatusCode(code int) bool {
	if code < 200 || code >= 400 {
		return false
	}
	return true
}

func (us *Uspacy) doRaw(url, method string, headers map[string]string, body io.Reader) ([]byte, error) {

	var (
		res *http.Response
		err error
	)

	if len(us.refreshToken) == 0 {
		us.isExpired = true
		us.refreshToken = us.bearerToken
		if us.TokenRefresh() == nil {
			return us.doRaw(url, method, headers, body)
		} else {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	switch us.isExpired {
	case true:
		req.Header.Add("Authorization", tokenPrefix+us.refreshToken)
	default:
		req.Header.Add("Authorization", tokenPrefix+us.bearerToken)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	for attempts := us.retries; attempts > 0; attempts-- {
		res, err = us.client.Do(req)
		if err == nil {
			break
		}

	}
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !handleStatusCode(res.StatusCode) {
		if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
			if !us.isExpired {
				us.isExpired = true
				if us.TokenRefresh() == nil {
					return us.doRaw(url, method, headers, body)
				} else {
					return nil, err
				}
			}
		}

		//errMsg := fmt.Sprintf("error occured while trying to (%s)\nbody - %s\ncode - %v\n", req.URL.String(), string(responseBody), res.StatusCode)
		//log.Println(errMsg)

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

func (us *Uspacy) doPostEncodedForm(url string, values url.Values) ([]byte, error) {
	var head = make(map[string]string)
	head["Content-Type"] = "application/x-www-form-urlencoded"
	head["Accept"] = "application/json"
	return us.doRaw(url, http.MethodPost, head, strings.NewReader(values.Encode()))
}

func (us *Uspacy) doDeleteEmptyHeaders(url string) (err error) {
	_, err = us.doRaw(url, http.MethodDelete, headersMap, nil)
	return err
}

func (us *Uspacy) buildURL(version, route string) string {
	return fmt.Sprintf("%s/%s/%s", us.mainHost, version, route)
}

func (us *Uspacy) doPostFormData(url string, textParams map[string]string, files map[string]io.ReadCloser) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	headers := make(map[string]string)

	for filename, file := range files {

		fileField, err := writer.CreateFormFile("files", filename)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(fileField, file)
		if err != nil {
			return nil, err
		}
	}

	for name, value := range textParams {
		writer.WriteField(name, value)
	}

	headers["Content-Type"] = writer.FormDataContentType()
	headers["Accept"] = "application/json"

	return us.doRaw(url, http.MethodPost, headers, body)

}
