package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/rand"
	"mime/multipart"
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
	lastStatusCode int
}

// errorLog stores unique error message and its attempts
type errorLog struct {
	message  string
	attempts []int
}

const (
	defaultClientTimeout = 20 * time.Second
	defaultRetries       = 3
	tokenPrefix          = "Bearer "
	defaultTimeout       = 5 * time.Second
	uploadTimeout        = 30 * time.Second

	// Backoff intervals for each retry attempt
	firstRetryDelay  = 5 * time.Second
	secondRetryDelay = 60 * time.Second
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
		isExpired: false,
	}
}

// prepareRequest creates and configures an HTTP request with appropriate headers and authorization
func (us *Uspacy) prepareRequest(ctx context.Context, url, method string, headers map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	// Set authorization token
	token := us.bearerToken
	if us.isExpired {
		token = us.RefreshToken
	}
	req.Header.Add("Authorization", tokenPrefix+token)

	// Add custom headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req.WithContext(ctx), nil
}

// calculateBackoff returns a fixed delay with jitter based on the attempt number
func (us *Uspacy) calculateBackoff(attempt int) time.Duration {
	var baseDelay time.Duration
	switch attempt {
	case 0: // First retry
		baseDelay = firstRetryDelay
	case 1: // Second retry
		baseDelay = secondRetryDelay
	default:
		baseDelay = secondRetryDelay
	}

	// Add jitter: random value between 75% and 100% of base delay
	jitterRange := baseDelay / 4
	jitter := time.Duration(rand.Int63n(int64(jitterRange)))
	return baseDelay - jitterRange + jitter
}

// executeRequest performs an HTTP request with retry mechanism
func (us *Uspacy) executeRequest(req *http.Request) ([]byte, int, error) {
	var (
		responseBody []byte
		errorLogs    = make(map[string]*errorLog)
	)

	for attempt := 0; attempt < defaultRetries; attempt++ {
		startTime := time.Now()

		res, err := us.client.Do(req)
		if err != nil {
			errMsg := err.Error()
			if log, exists := errorLogs[errMsg]; exists {
				log.attempts = append(log.attempts, attempt+1)
			} else {
				errorLogs[errMsg] = &errorLog{
					message:  errMsg,
					attempts: []int{attempt + 1},
				}
			}

			backoff := us.calculateBackoff(attempt)
			if time.Since(startTime) > backoff+5*time.Second {
				break
			}

			time.Sleep(backoff)
			continue
		}

		defer res.Body.Close()
		responseBody, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, 0, err
		}

		us.lastStatusCode = res.StatusCode
		return responseBody, res.StatusCode, nil
	}

	// Format unique errors with their attempts
	var errorDetails strings.Builder
	for _, log := range errorLogs {
		errorDetails.WriteString(fmt.Sprintf("Error '%s' on attempts: %v\n",
			log.message, log.attempts))
	}

	return nil, 0, fmt.Errorf("request failed after %d retries:\n%s",
		defaultRetries, errorDetails.String())
}

// doRaw performs an HTTP request with token handling and retry mechanism
func (us *Uspacy) doRaw(url, method string, headers map[string]string, body io.Reader, timeout time.Duration) ([]byte, int, error) {
	if len(us.RefreshToken) == 0 {
		us.isExpired = true
		us.RefreshToken = us.bearerToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := us.prepareRequest(ctx, url, method, headers, body)
	if err != nil {
		return nil, 0, err
	}

	responseBody, statusCode, err := us.executeRequest(req)
	if err != nil {
		return nil, 0, err
	}

	if statusCode == http.StatusUnauthorized && !us.isExpired {
		if _, err := us.TokenRefresh(); err != nil {
			return nil, statusCode, err
		}
		// Retry request with new token
		req, err = us.prepareRequest(ctx, url, method, headers, body)
		if err != nil {
			return nil, 0, err
		}
		responseBody, statusCode, err = us.executeRequest(req)
		if err != nil {
			return nil, statusCode, err
		}
	}

	if statusCode < 200 || statusCode >= 400 {
		return responseBody, statusCode, errors.New(string(responseBody))
	}

	us.isExpired = false
	return responseBody, statusCode, nil
}

// doGetEmptyHeaders performs a GET request with default headers
func (us *Uspacy) doGetEmptyHeaders(url string) ([]byte, error) {
	// Perform a GET request with default headers
	response, _, err := us.doRaw(url, http.MethodGet, headersMap, nil, defaultTimeout)
	return response, err
}

// doPost performs a POST request with JSON body and optional additional headers
func (us *Uspacy) doPost(url string, body interface{}, headers ...map[string]string) ([]byte, int, error) {
	// Encode the JSON body
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Merge default headers with additional headers
	requestHeaders := make(map[string]string)
	maps.Copy(requestHeaders, headersMap)

	for _, headerMap := range headers {
		for key, value := range headerMap {
			if value != "" {
				requestHeaders[key] = value
			}
		}
	}

	// Perform the POST request
	response, code, err := us.doRaw(url, http.MethodPost, requestHeaders, &buf, defaultTimeout)
	return response, code, err
}

// doPatchEmptyHeaders performs a PATCH request with default headers and JSON body
func (us *Uspacy) doPatchEmptyHeaders(url string, body interface{}) ([]byte, error) {
	// Encode the JSON body
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}

	// Perform the PATCH request
	response, _, err := us.doRaw(url, http.MethodPatch, headersMap, &buf, defaultTimeout)
	return response, err
}

// doPostEncodedForm performs a POST request with form-encoded data
func (us *Uspacy) doPostEncodedForm(url string, values url.Values) ([]byte, error) {
	// Set the Content-Type header to application/x-www-form-urlencoded
	var head = make(map[string]string)
	head["Content-Type"] = "application/x-www-form-urlencoded"
	head["Accept"] = "application/json"

	// Perform the POST request
	response, _, err := us.doRaw(url, http.MethodPost, head, strings.NewReader(values.Encode()), defaultTimeout)
	return response, err
}

// doDeleteEmptyHeaders performs a DELETE request with default headers and optional JSON body
func (us *Uspacy) doDeleteEmptyHeaders(url string, body interface{}) (int, error) {
	// Encode the JSON body if provided
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	// Perform the DELETE request
	_, code, err := us.doRaw(url, http.MethodDelete, headersMap, nil, defaultTimeout)
	return code, err
}

// doPostFormData performs a multipart form POST request with files and text parameters.
// Returns error if no files are provided or if all provided files are invalid.
// Uses extended timeout for handling large file uploads.
func (us *Uspacy) doPostFormData(url string, textParams map[string]string, files map[string]io.ReadCloser) ([]byte, error) {
	// Check if files are provided
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided for upload")
	}

	// Create a multipart form writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	headers := make(map[string]string)

	// Add files to the form
	validFilesCount := 0
	for filename, file := range files {
		if file == nil || filename == "" {
			continue
		}
		fileField, err := writer.CreateFormFile("files[]", filename)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(fileField, file); err != nil {
			return nil, err
		}
		validFilesCount++
	}

	// Check if any valid files were added
	if validFilesCount == 0 {
		return nil, fmt.Errorf("no valid files found for upload")
	}

	// Add text parameters to the form
	for key, value := range textParams {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	// Set the Content-Type header to the form's content type
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Accept"] = "application/json"

	// Perform the POST request with extended timeout
	response, _, err := us.doRaw(url, http.MethodPost, headers, body, uploadTimeout)
	return response, err
}

// buildURL constructs a full URL from version and route components
func (us *Uspacy) buildURL(version, route string) string {
	return fmt.Sprintf("%s/%s/%s", us.mainHost, version, route)
}
