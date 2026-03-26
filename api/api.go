package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Uspacy struct {
	bearerToken    string
	RefreshToken   string
	client         *http.Client
	mainHost       string
	lastStatusCode int
	mu             sync.RWMutex
}

// errorLog stores unique error message and its attempts
type errorLog struct {
	message  string
	attempts []int
}

const (
	defaultClientTimeout = 30 * time.Second
	defaultRetries       = 3
	tokenPrefix          = "Bearer "
	maxRetryAfter        = 5 * time.Minute

	// Backoff intervals for each retry attempt
	firstRetryDelay  = 3 * time.Second
	secondRetryDelay = 5 * time.Second
)

// New creates an Uspacy object
func New(token, refresh, host string) *Uspacy {
	bearerToken := strings.TrimPrefix(token, tokenPrefix)
	refreshToken := strings.TrimPrefix(refresh, tokenPrefix)

	// Fallback: if refresh token is empty, use bearer token
	if len(refreshToken) == 0 {
		refreshToken = bearerToken
	}

	return &Uspacy{
		bearerToken:  bearerToken,
		RefreshToken: refreshToken,
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
		mainHost: host,
	}
}

// prepareRequest creates and configures an HTTP request with appropriate headers and authorization
func (us *Uspacy) prepareRequest(url, method string, headers map[string]string, body []byte) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if len(headers) == 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	// Set authorization token
	us.mu.RLock()
	token := us.bearerToken
	us.mu.RUnlock()
	req.Header.Add("Authorization", tokenPrefix+token)

	// Add custom headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
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

// doRaw performs an HTTP request with token handling and retry mechanism
func (us *Uspacy) doRaw(url, method string, headers map[string]string, body []byte) ([]byte, int, error) {
	return us.doRawInternal(url, method, headers, body, false)
}

// doRawSkipRefresh performs an HTTP request without token refresh on 401
func (us *Uspacy) doRawSkipRefresh(url, method string, headers map[string]string, body []byte) ([]byte, int, error) {
	return us.doRawInternal(url, method, headers, body, true)
}

// doRawInternal performs an HTTP request with optional token refresh
func (us *Uspacy) doRawInternal(url, method string, headers map[string]string, body []byte, skipTokenRefresh bool) ([]byte, int, error) {
	var (
		responseBody   []byte
		statusCode     int
		errorLogs      = make(map[string]*errorLog)
		tokenRefreshed = false
	)

	for attempt := 0; attempt < defaultRetries; attempt++ {
		req, err := us.prepareRequest(url, method, headers, body)
		if err != nil {
			return nil, 0, err
		}

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

			if attempt < defaultRetries-1 {
				time.Sleep(us.calculateBackoff(attempt))
			}
			continue
		}

		responseBody, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, 0, err
		}

		statusCode = res.StatusCode
		us.lastStatusCode = statusCode

		// Handle 401 Unauthorized - refresh token and retry (only once)
		if statusCode == http.StatusUnauthorized && !skipTokenRefresh && !tokenRefreshed {
			if _, err := us.TokenRefresh(); err != nil {
				return nil, statusCode, err
			}
			tokenRefreshed = true
			attempt-- // Don't consume retry attempt for token refresh
			continue
		}

		// Handle 429 Too Many Requests - retry with backoff
		if statusCode == http.StatusTooManyRequests {
			retryAfter := us.parseRetryAfter(res.Header.Get("Retry-After"))
			if attempt < defaultRetries-1 {
				time.Sleep(retryAfter)
			}
			continue
		}

		// Success or non-retryable error
		if statusCode >= 200 && statusCode < 500 {
			break
		}

		// Server error (5xx) - retry
		if attempt < defaultRetries-1 {
			time.Sleep(us.calculateBackoff(attempt))
		}
	}

	// Check if all retries failed with errors
	if len(errorLogs) > 0 && statusCode == 0 {
		var errorDetails strings.Builder
		for _, log := range errorLogs {
			fmt.Fprintf(&errorDetails, "Error '%s' on attempts: %v\n",
				log.message, log.attempts)
		}
		return nil, 0, fmt.Errorf("request failed after %d retries:\n%s",
			defaultRetries, errorDetails.String())
	}

	if statusCode < 200 || statusCode >= 400 {
		return responseBody, statusCode, fmt.Errorf("request failed: [%s] %s, status code: %d, response: %s", method, url, statusCode, string(responseBody))
	}

	return responseBody, statusCode, nil
}

// parseRetryAfter parses Retry-After header value and returns duration
func (us *Uspacy) parseRetryAfter(header string) time.Duration {
	if header == "" {
		return us.calculateBackoff(0)
	}

	var duration time.Duration

	// Try parsing as seconds
	if seconds, err := time.ParseDuration(header + "s"); err == nil {
		duration = seconds
	} else if t, err := http.ParseTime(header); err == nil {
		// Try parsing as HTTP-date
		duration = time.Until(t)
		if duration <= 0 {
			return us.calculateBackoff(0)
		}
	} else {
		return us.calculateBackoff(0)
	}

	// Cap the duration at maxRetryAfter
	if duration > maxRetryAfter {
		return maxRetryAfter
	}

	return duration
}

// doGetEmptyHeaders performs a GET request with default headers and optional additional headers
func (us *Uspacy) doGetEmptyHeaders(url string, headers ...map[string]string) ([]byte, error) {
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

	response, _, err := us.doRaw(url, http.MethodGet, requestHeaders, nil)
	return response, err
}

// doPost performs a POST request with JSON body and optional additional headers
func (us *Uspacy) doPost(url string, body any, headers ...map[string]string) ([]byte, int, error) {
	jsonBody, err := json.Marshal(body)
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

	response, code, err := us.doRaw(url, http.MethodPost, requestHeaders, jsonBody)
	return response, code, err
}

// doPatchEmptyHeaders performs a PATCH request with default headers and JSON body
func (us *Uspacy) doPatchEmptyHeaders(url string, body any) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	response, _, err := us.doRaw(url, http.MethodPatch, headersMap, jsonBody)
	return response, err
}

// doPostEncodedForm performs a POST request with form-encoded data
func (us *Uspacy) doPostEncodedForm(url string, values url.Values) ([]byte, error) {
	head := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	}

	response, _, err := us.doRaw(url, http.MethodPost, head, []byte(values.Encode()))
	return response, err
}

// doDeleteEmptyHeaders performs a DELETE request with default headers and optional JSON body
func (us *Uspacy) doDeleteEmptyHeaders(url string, body any) (int, error) {
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	_, code, err := us.doRaw(url, http.MethodDelete, headersMap, jsonBody)
	return code, err
}

// doPostFormData performs a multipart form POST request with files and text parameters.
// Returns error if no files are provided or if all provided files are invalid.
func (us *Uspacy) doPostFormData(url string, textParams map[string]string, files map[string]io.ReadCloser) ([]byte, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided for upload")
	}

	// Ensure all files are closed when done
	defer func() {
		for _, file := range files {
			if file != nil {
				file.Close()
			}
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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

	if validFilesCount == 0 {
		return nil, fmt.Errorf("no valid files found for upload")
	}

	// Add text parameters to the form
	for key, value := range textParams {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Accept":       "application/json",
	}

	response, _, err := us.doRaw(url, http.MethodPost, headers, body.Bytes())
	return response, err
}

// buildURL constructs a full URL by joining all parts with "/"
func (us *Uspacy) buildURL(parts ...string) string {
	allParts := append([]string{us.mainHost}, parts...)
	var result strings.Builder

	for i, part := range allParts {
		// Trim slashes from both ends of the part
		trimmed := strings.Trim(part, "/")
		if trimmed == "" {
			continue
		}

		if i > 0 {
			result.WriteString("/")
		}
		result.WriteString(trimmed)
	}

	return result.String()
}
