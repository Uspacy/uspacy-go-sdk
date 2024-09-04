package api

import (
	"bytes"
	"context"
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
	bearerToken    string
	RefreshToken   string
	client         *http.Client
	mainHost       string
	isExpired      bool
	retries        int `default0:"3"`
	lastStatusCode int
}

const (
	defaultClientTimeout = 20 * time.Second
	defaultRetries       = 3
	tokenPrefix          = "Bearer "
)

// New creates an Uspacy object
func New(token, refresh, host string) *Uspacy {

	return &Uspacy{
		bearerToken:  strings.TrimPrefix(token, tokenPrefix),
		RefreshToken: strings.TrimPrefix(refresh, tokenPrefix),
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

func (us *Uspacy) doRaw(url, method string, headers map[string]string, body io.Reader) ([]byte, int, error) {

	var (
		res          *http.Response
		err          error
		responseBody []byte
		errorDetails strings.Builder
	)

	if len(us.RefreshToken) == 0 {
		us.isExpired = true
		us.RefreshToken = us.bearerToken
		if _, errRefresh := us.TokenRefresh(); errRefresh == nil {
			return us.doRaw(url, method, headers, body)
		} else {
			return nil, 0, err
		}
	}

	// Create a context with a cancel function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	if len(headers) == 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	switch us.isExpired {
	case true:
		req.Header.Add("Authorization", tokenPrefix+us.RefreshToken)
	default:
		req.Header.Add("Authorization", tokenPrefix+us.bearerToken)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	req = req.WithContext(ctx)

	backoff := 3 * time.Second

	for attempts := 3; attempts > 0; attempts-- {
		// Start a new context with timeout for this attempt
		attemptCtx, attemptCancel := context.WithTimeout(ctx, backoff+5*time.Second)
		defer attemptCancel()

		// Create a new request with the attempt context
		attemptReq := req.WithContext(attemptCtx)

		startTime := time.Now()

		res, err = us.client.Do(attemptReq)

		if err == nil {
			defer res.Body.Close()
			responseBody, err = io.ReadAll(res.Body)
			if err != nil {
				return nil, 0, err
			}
			// Set last status code
			us.lastStatusCode = res.StatusCode
			break
		}

		// Record error details
		errorDetails.WriteString(fmt.Sprintf("Attempt %d: %s\n", us.retries-attempts+1, err.Error()))

		// Check if the attempt took too long
		if time.Since(startTime) > backoff+5*time.Second {
			break
		}

		// Progressive delay (exponential backoff)
		time.Sleep(backoff)
		backoff *= 2
	}

	if err != nil {
		return nil, 0, fmt.Errorf("request failed after %d retries: %s", us.retries, errorDetails.String())
	}

	if !handleStatusCode(res.StatusCode) {
		if res.StatusCode == http.StatusUnauthorized {
			if !us.isExpired {
				us.isExpired = true
				if _, errRefresh := us.TokenRefresh(); errRefresh == nil {
					return us.doRaw(url, method, headers, body)
				} else {
					return nil, 0, err
				}
			}
		}
		return responseBody, us.lastStatusCode, errors.New(string(responseBody))
	}
	us.isExpired = false
	return responseBody, us.lastStatusCode, nil
}

func (us *Uspacy) doGetEmptyHeaders(url string) ([]byte, error) {
	response, _, err := us.doRaw(url, http.MethodGet, headersMap, nil)
	return response, err
}

func (us *Uspacy) doPostEmptyHeaders(url string, body interface{}) ([]byte, int, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, 0, err
	}
	response, code, err := us.doRaw(url, http.MethodPost, headersMap, &buf)
	return response, code, err
}

func (us *Uspacy) doPatchEmptyHeaders(url string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}
	response, _, err := us.doRaw(url, http.MethodPatch, headersMap, &buf)
	return response, err
}

func (us *Uspacy) doPostEncodedForm(url string, values url.Values) ([]byte, error) {
	var head = make(map[string]string)
	head["Content-Type"] = "application/x-www-form-urlencoded"
	head["Accept"] = "application/json"
	response, _, err := us.doRaw(url, http.MethodPost, head, strings.NewReader(values.Encode()))
	return response, err
}

func (us *Uspacy) doDeleteEmptyHeaders(url string) (int, error) {
	_, code, err := us.doRaw(url, http.MethodDelete, headersMap, nil)
	return code, err
}

func (us *Uspacy) buildURL(version, route string) string {
	return fmt.Sprintf("%s/%s/%s", us.mainHost, version, route)
}

func (us *Uspacy) doPostFormData(url string, textParams map[string]string, files map[string]io.ReadCloser) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	headers := make(map[string]string)
	for filename, file := range files {

		fileField, err := writer.CreateFormFile("files[]", filename)
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
	response, _, err := us.doRaw(url, http.MethodPost, headers, body)
	return response, err
}
