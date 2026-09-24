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

// prepareRequestCtx creates and configures an HTTP request with appropriate headers and
// authorization, bound to ctx so the caller can cancel it or let it time out.
func (us *Uspacy) prepareRequestCtx(ctx context.Context, url, method string, headers map[string]string, body []byte) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
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

// prepareRequest creates and configures an HTTP request with appropriate headers and authorization.
// It delegates to prepareRequestCtx with context.Background(), so its behaviour is unchanged.
func (us *Uspacy) prepareRequest(url, method string, headers map[string]string, body []byte) (*http.Request, error) {
	return us.prepareRequestCtx(context.Background(), url, method, headers, body)
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

// doRawInternal performs an HTTP request with optional token refresh.
// It delegates to doRawInternalCtx with context.Background(), so its behaviour is unchanged.
func (us *Uspacy) doRawInternal(url, method string, headers map[string]string, body []byte, skipTokenRefresh bool) ([]byte, int, error) {
	return us.doRawInternalCtx(context.Background(), url, method, headers, body, skipTokenRefresh)
}

// doRawInternalCtx performs an HTTP request with optional token refresh. It carries ctx on
// every request it sends (including the token refresh), and every backoff or Retry-After
// wait between attempts stops as soon as ctx is done, instead of sleeping in full.
func (us *Uspacy) doRawInternalCtx(ctx context.Context, url, method string, headers map[string]string, body []byte, skipTokenRefresh bool) ([]byte, int, error) {
	var (
		responseBody   []byte
		statusCode     int
		errorLogs      = make(map[string]*errorLog)
		tokenRefreshed = false
		lastErr        error // most recent attempt's failure: *HTTPError, or a transport error
	)

	for attempt := 0; attempt < defaultRetries; attempt++ {
		req, err := us.prepareRequestCtx(ctx, url, method, headers, body)
		if err != nil {
			return nil, 0, err
		}

		res, err := us.client.Do(req)
		if err != nil {
			lastErr = err
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
				if sleepErr := sleepCtx(ctx, us.calculateBackoff(attempt)); sleepErr != nil {
					return nil, statusCode, abortedWhileWaiting(lastErr, sleepErr)
				}
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
			if _, err := us.tokenRefreshCtx(ctx); err != nil {
				return nil, statusCode, err
			}
			tokenRefreshed = true
			attempt-- // Don't consume retry attempt for token refresh
			continue
		}

		// Handle 429 Too Many Requests - retry with backoff
		if statusCode == http.StatusTooManyRequests {
			lastErr = &HTTPError{Method: method, URL: url, StatusCode: statusCode, Body: responseBody}
			retryAfter := us.parseRetryAfter(res.Header.Get("Retry-After"))
			if attempt < defaultRetries-1 {
				if sleepErr := sleepCtx(ctx, retryAfter); sleepErr != nil {
					return nil, statusCode, abortedWhileWaiting(lastErr, sleepErr)
				}
			}
			continue
		}

		// Success or non-retryable error
		if statusCode >= 200 && statusCode < 500 {
			break
		}

		// Server error (5xx) - retry
		lastErr = &HTTPError{Method: method, URL: url, StatusCode: statusCode, Body: responseBody}
		if attempt < defaultRetries-1 {
			if sleepErr := sleepCtx(ctx, us.calculateBackoff(attempt)); sleepErr != nil {
				return nil, statusCode, abortedWhileWaiting(lastErr, sleepErr)
			}
		}
	}

	// Check if all retries failed with errors
	if len(errorLogs) > 0 && statusCode == 0 {
		return nil, 0, requestFailedAfterRetries(ctx, defaultRetries, errorLogs)
	}

	if statusCode < 200 || statusCode >= 400 {
		return responseBody, statusCode, &HTTPError{Method: method, URL: url, StatusCode: statusCode, Body: responseBody}
	}

	return responseBody, statusCode, nil
}

// sleepCtx waits d, or returns ctx.Err() as soon as ctx is done.
func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// abortedWhileWaiting builds the error doRawInternalCtx returns when ctx ends while it is
// sleeping before a retry. lastErr is the most recent attempt's failure: if it is already
// an *HTTPError (the previous attempt got a 429 or 5xx response), that response is still
// meaningful on its own and is returned as is. Otherwise lastErr is a transport-level error
// (e.g. a connection failure) with no HTTP response behind it, so ctxErr is returned with
// lastErr folded into its text, keeping errors.Is(err, ctx.Err()) working through the %w chain.
func abortedWhileWaiting(lastErr error, ctxErr error) error {
	var httpErr *HTTPError
	if errors.As(lastErr, &httpErr) {
		return httpErr
	}
	if lastErr != nil {
		return fmt.Errorf("request aborted while waiting to retry after %v: %w", lastErr, ctxErr)
	}
	return ctxErr
}

// requestFailedAfterRetries builds the error doRawInternalCtx returns when every attempt
// failed at the transport level, so no HTTP response was ever received (statusCode stays
// 0). The last attempt runs no backoff sleep, so a context that ends during it is never
// caught by sleepCtx/abortedWhileWaiting; if ctx has already ended by the time this runs,
// ctx.Err() is folded into the text with %w, keeping errors.Is(err, ctx.Err()) working.
// Every existing caller runs this under context.Background(), which never ends, so for
// them the text and type are exactly what doRawInternalCtx returned before ctx support
// was added.
func requestFailedAfterRetries(ctx context.Context, retries int, errorLogs map[string]*errorLog) error {
	var errorDetails strings.Builder
	for _, log := range errorLogs {
		fmt.Fprintf(&errorDetails, "Error '%s' on attempts: %v\n",
			log.message, log.attempts)
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return fmt.Errorf("request failed after %d retries:\n%scontext error: %w",
			retries, errorDetails.String(), ctxErr)
	}
	return fmt.Errorf("request failed after %d retries:\n%s",
		retries, errorDetails.String())
}

// HTTPError is returned for a final non-2xx response from doRawInternalCtx. Error() keeps
// the exact historical text so existing string parsing keeps working, even though the error
// is now a concrete type: "request failed: [GET] <url>, status code: 403, response: <body>".
type HTTPError struct {
	Method     string
	URL        string
	StatusCode int
	Body       []byte
	Err        error // optional cause, e.g. a failed token refresh after a 401
}

// Error returns text identical to the SDK's historical inline error message. It does not
// include Err, so wrapping a cause (see asHTTPError) never changes this text.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("request failed: [%s] %s, status code: %d, response: %s", e.Method, e.URL, e.StatusCode, string(e.Body))
}

// Unwrap exposes Err so callers can use errors.Is / errors.As on the cause.
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// asHTTPError normalizes a (method, url, statusCode, err) result so every non-2xx failure
// is an *HTTPError carrying Method and URL. If err is nil, or already an *HTTPError (or
// wraps one), it is returned unchanged. Otherwise, for statusCode >= 400, it is wrapped as
// &HTTPError{Method: method, URL: url, StatusCode: statusCode, Err: err} — for example the
// raw error doRawInternalCtx returns when a 401's token refresh fails. Unexported helper
// for the context-aware CRM methods (GetFieldsCtx, GetEntityCtx) that need a single error
// type.
func asHTTPError(method, url string, statusCode int, err error) error {
	if err == nil {
		return nil
	}
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return err
	}
	if statusCode >= 400 {
		return &HTTPError{Method: method, URL: url, StatusCode: statusCode, Err: err}
	}
	return err
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
