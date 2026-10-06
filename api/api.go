package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Uspacy struct {
	bearerToken   string
	refreshToken  string
	client        *http.Client
	mainHost      string
	mu            sync.RWMutex
	maxRetries    int
	retryBase     time.Duration
	retryMax      time.Duration
	refreshFlight singleflight.Group
}

// errorLog stores unique error message and its attempts
type errorLog struct {
	message  string
	attempts []int
}

// Option configures an Uspacy client.
type Option func(*Uspacy)

// RequestOption configures an individual HTTP request.
type RequestOption func(*requestConfig)

type requestConfig struct {
	headers map[string]string
}

// WithHeader adds a single HTTP header to the request. An empty value is ignored.
func WithHeader(key, value string) RequestOption {
	return func(cfg *requestConfig) {
		if value == "" {
			return
		}
		if cfg.headers == nil {
			cfg.headers = make(map[string]string)
		}
		cfg.headers[key] = value
	}
}

// WithHeaders adds all headers from h to the request. Empty values are ignored.
func WithHeaders(h map[string]string) RequestOption {
	return func(cfg *requestConfig) {
		for k, v := range h {
			WithHeader(k, v)(cfg)
		}
	}
}

func newRequestConfig(opts ...RequestOption) requestConfig {
	cfg := requestConfig{headers: make(map[string]string)}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

const (
	defaultClientTimeout = 30 * time.Second
	defaultMaxRetries    = 3
	tokenPrefix          = "Bearer "
	maxRetryAfter        = 5 * time.Minute
	defaultRetryBase     = 3 * time.Second
	defaultRetryMax      = 30 * time.Second
)

// New creates an Uspacy object. Options can be passed to customize retry behavior.
func New(token, refresh, host string, opts ...Option) *Uspacy {
	bearerToken := strings.TrimPrefix(token, tokenPrefix)
	refreshToken := strings.TrimPrefix(refresh, tokenPrefix)

	// Fallback: if refresh token is empty, use bearer token
	if len(refreshToken) == 0 {
		refreshToken = bearerToken
	}

	us := &Uspacy{
		bearerToken:  bearerToken,
		refreshToken: refreshToken,
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
		mainHost:   host,
		maxRetries: defaultMaxRetries,
		retryBase:  defaultRetryBase,
		retryMax:   defaultRetryMax,
	}

	for _, opt := range opts {
		opt(us)
	}

	return us
}

// RefreshToken returns the current refresh token. Safe for concurrent use.
func (us *Uspacy) RefreshToken() string {
	us.mu.RLock()
	defer us.mu.RUnlock()
	return us.refreshToken
}

// WithMaxRetries sets the maximum number of request attempts (must be > 0).
func WithMaxRetries(n int) Option {
	return func(us *Uspacy) {
		if n > 0 {
			us.maxRetries = n
		}
	}
}

// WithRetryBackoff sets the base and maximum delays between retries. The actual delay
// for attempt i is base * 2^i with equal jitter, capped at max.
func WithRetryBackoff(base, max time.Duration) Option {
	return func(us *Uspacy) {
		if base > 0 {
			us.retryBase = base
		}
		if max > 0 {
			us.retryMax = max
		}
	}
}

// WithHTTPClient sets the HTTP client used for requests. By default a client with a
// 30-second timeout is created.
func WithHTTPClient(c *http.Client) Option {
	return func(us *Uspacy) {
		if c != nil {
			us.client = c
		}
	}
}

// prepareRequest creates and configures an HTTP request with appropriate headers and
// authorization, bound to ctx so the caller can cancel it or let it time out.
func (us *Uspacy) prepareRequest(ctx context.Context, url, method string, headers map[string]string, body []byte) (*http.Request, error) {
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

// responseClass tells doRequest what to do with a single HTTP response.
type responseClass int

const (
	responseDone responseClass = iota
	responseRefresh
	responseRetry
)

// doRaw performs an HTTP request with token handling and retry mechanism. Every final
// non-2xx answer is returned as *HTTPError, including 3xx and a 401 whose token refresh
// failed. If ctx ends before the request completes, the error preserves the last 429/5xx
// response when available and wraps ctx.Err() so errors.Is(err, ctx.Err()) works.
func (us *Uspacy) doRaw(ctx context.Context, url, method string, headers map[string]string, body []byte) ([]byte, int, error) {
	return us.doRequest(ctx, url, method, headers, body, false)
}

// doRequest performs an HTTP request with optional token refresh. It carries ctx on every
// request it sends (including the token refresh), and every backoff or Retry-After wait
// between attempts stops as soon as ctx is done. A 401 whose token refresh fails is
// reported as an *HTTPError with the refresh error as its cause. A final 3xx is also an
// *HTTPError. If ctx ends, the returned error wraps ctx.Err() and, when a 429/5xx
// response was already seen, carries that response so callers can inspect it.
func (us *Uspacy) doRequest(ctx context.Context, url, method string, headers map[string]string, body []byte, skipTokenRefresh bool) ([]byte, int, error) {
	var (
		responseBody   []byte
		statusCode     int
		errorLogs      = make(map[string]*errorLog)
		tokenRefreshed = false
		lastHTTPErr    *HTTPError // last 429 or 5xx response seen; a 401 a token refresh resolved never sets this
	)

	for attempt := 0; attempt < us.maxRetries; attempt++ {
		req, err := us.prepareRequest(ctx, url, method, headers, body)
		if err != nil {
			return nil, 0, err
		}

		res, err := us.client.Do(req)
		if err != nil {
			us.logError(errorLogs, err.Error(), attempt+1)
			if us.shouldRetry(attempt) {
				if sleepErr := sleep(ctx, us.nextBackoff(attempt, 0, false)); sleepErr != nil {
					return us.abort(method, url, lastHTTPErr, sleepErr)
				}
			} else if ctxErr := ctx.Err(); ctxErr != nil {
				return us.abort(method, url, lastHTTPErr, ctxErr)
			}
			continue
		}

		responseBody, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return us.abort(method, url, lastHTTPErr, ctxErr)
			}
			return nil, 0, err
		}

		statusCode = res.StatusCode

		switch us.classifyResponse(method, statusCode, skipTokenRefresh, tokenRefreshed) {
		case responseRefresh:
			if _, err := us.tokenRefresh(ctx); err != nil {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return us.abort(method, url, lastHTTPErr, ctxErr)
				}
				return nil, statusCode, &HTTPError{Method: method, URL: url, StatusCode: http.StatusUnauthorized, Err: err}
			}
			tokenRefreshed = true
			attempt-- // Don't consume retry attempt for token refresh
			continue
		case responseRetry:
			lastHTTPErr = &HTTPError{Method: method, URL: url, StatusCode: statusCode, Body: responseBody}
			if us.shouldRetry(attempt) {
				retryAfter, retryAfterSet := time.Duration(0), false
				if statusCode == http.StatusTooManyRequests {
					retryAfter, retryAfterSet = us.parseRetryAfter(res.Header.Get("Retry-After"))
				}
				if sleepErr := sleep(ctx, us.nextBackoff(attempt, retryAfter, retryAfterSet)); sleepErr != nil {
					return us.abort(method, url, lastHTTPErr, sleepErr)
				}
			}
			continue
		case responseDone:
			return us.finalResult(method, url, statusCode, responseBody, errorLogs, lastHTTPErr, ctx)
		}
	}

	return us.finalResult(method, url, statusCode, responseBody, errorLogs, lastHTTPErr, ctx)
}

// classifyResponse decides what to do with a single HTTP response.
func (us *Uspacy) classifyResponse(method string, statusCode int, skipTokenRefresh, tokenRefreshed bool) responseClass {
	switch {
	case statusCode == http.StatusUnauthorized && !skipTokenRefresh && !tokenRefreshed:
		return responseRefresh
	case statusCode == http.StatusTooManyRequests:
		return responseRetry
	case statusCode >= 500 && isIdempotent(method):
		return responseRetry
	default:
		return responseDone
	}
}

// isIdempotent reports whether method is safe to retry after a server error.
func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

// shouldRetry reports whether there are still attempts left after the current one.
func (us *Uspacy) shouldRetry(attempt int) bool {
	return attempt < us.maxRetries-1
}

// logError records a transport-level error against its attempt number.
func (us *Uspacy) logError(errorLogs map[string]*errorLog, message string, attempt int) {
	if log, exists := errorLogs[message]; exists {
		log.attempts = append(log.attempts, attempt)
	} else {
		errorLogs[message] = &errorLog{
			message:  message,
			attempts: []int{attempt},
		}
	}
}

// sleep waits d, or returns ctx.Err() as soon as ctx is done. An already-done ctx wins
// even when d is not positive (e.g. Retry-After: 0), so no further attempt starts after ctx
// has ended; otherwise a non-positive d returns at once, without a timer.
func sleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// nextBackoff returns the wait duration before the next attempt. If retryAfterSet is
// true the header value is used directly (capped at maxRetryAfter); otherwise the delay
// is exponential (base * 2^attempt) with equal jitter, capped at retryMax.
func (us *Uspacy) nextBackoff(attempt int, retryAfter time.Duration, retryAfterSet bool) time.Duration {
	if retryAfterSet {
		if retryAfter > maxRetryAfter {
			return maxRetryAfter
		}
		return retryAfter
	}

	var shift uint
	if attempt > 30 {
		shift = 30
	} else {
		shift = uint(attempt)
	}
	delay := us.retryBase * (1 << shift)
	if delay > us.retryMax || delay <= 0 {
		delay = us.retryMax
	}
	if delay <= 0 {
		return 0
	}

	// Equal jitter: [delay/2, delay).
	half := delay / 2
	return half + time.Duration(rand.Float64()*float64(half))
}

// abort reports ctx ending before a retry could happen. When a 429/5xx response was
// already seen, it is preserved as an *HTTPError with ctxErr attached via Err so both
// errors.As(*HTTPError) and errors.Is(ctxErr) work. Otherwise it returns a wrapped
// context error.
func (us *Uspacy) abort(method, url string, lastHTTPErr *HTTPError, ctxErr error) ([]byte, int, error) {
	if lastHTTPErr != nil {
		errWithCtx := *lastHTTPErr
		errWithCtx.Err = ctxErr
		return nil, lastHTTPErr.StatusCode, &errWithCtx
	}
	return nil, 0, fmt.Errorf("request aborted: %w", ctxErr)
}

// finalResult builds the return value once the retry loop is done.
func (us *Uspacy) finalResult(method, url string, statusCode int, body []byte, errorLogs map[string]*errorLog, lastHTTPErr *HTTPError, ctx context.Context) ([]byte, int, error) {
	if len(errorLogs) > 0 && statusCode == 0 {
		return nil, 0, requestFailedAfterRetries(ctx, us.maxRetries, errorLogs)
	}
	if statusCode < 200 || statusCode >= 300 {
		if lastHTTPErr != nil && lastHTTPErr.StatusCode == statusCode {
			return body, statusCode, lastHTTPErr
		}
		return body, statusCode, &HTTPError{Method: method, URL: url, StatusCode: statusCode, Body: body}
	}
	return body, statusCode, nil
}

// requestFailedAfterRetries builds the error when every attempt failed at the transport
// level, so no HTTP response was ever received.
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

// HTTPError represents a failed HTTP call. Every public method reports a final non-2xx
// answer as *HTTPError, including 3xx and a 401 whose token refresh failed (the refresh
// error is attached via Err). Error() keeps the historical text so existing string parsing
// keeps working. Use errors.As for *HTTPError and errors.Is on Err for context errors.
type HTTPError struct {
	Method     string
	URL        string
	StatusCode int
	Body       []byte
	Err        error // optional cause, e.g. a failed token refresh after a 401 or a context error
}

// Error returns text identical to the SDK's historical inline error message. It does not
// include Err, so wrapping a cause never changes this text.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("request failed: [%s] %s, status code: %d, response: %s", e.Method, e.URL, e.StatusCode, string(e.Body))
}

// Unwrap exposes Err so callers can use errors.Is / errors.As on the cause.
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// parseRetryAfter parses Retry-After header value. The bool reports whether the header
// was present and valid; when false, the caller should fall back to exponential backoff.
func (us *Uspacy) parseRetryAfter(header string) (time.Duration, bool) {
	if header == "" {
		return 0, false
	}

	// Try parsing as seconds.
	if seconds, err := time.ParseDuration(header + "s"); err == nil {
		if seconds > maxRetryAfter {
			return maxRetryAfter, true
		}
		return seconds, true
	}

	// Try parsing as HTTP-date.
	if t, err := http.ParseTime(header); err == nil {
		duration := time.Until(t)
		if duration <= 0 {
			return 0, false
		}
		if duration > maxRetryAfter {
			return maxRetryAfter, true
		}
		return duration, true
	}

	return 0, false
}

// doGetEmptyHeaders performs a GET request with default headers and optional additional headers.
func (us *Uspacy) doGetEmptyHeaders(ctx context.Context, url string, opts ...RequestOption) ([]byte, error) {
	requestHeaders := mergeHeaders(opts)
	response, _, err := us.doRaw(ctx, url, http.MethodGet, requestHeaders, nil)
	return response, err
}

// doPost performs a POST request with JSON body and optional additional headers.
func (us *Uspacy) doPost(ctx context.Context, url string, body any, opts ...RequestOption) ([]byte, int, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	requestHeaders := mergeHeaders(opts)
	response, code, err := us.doRaw(ctx, url, http.MethodPost, requestHeaders, jsonBody)
	return response, code, err
}

// doPatchEmptyHeaders performs a PATCH request with default headers and JSON body.
func (us *Uspacy) doPatchEmptyHeaders(ctx context.Context, url string, body any, opts ...RequestOption) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	requestHeaders := mergeHeaders(opts)
	response, _, err := us.doRaw(ctx, url, http.MethodPatch, requestHeaders, jsonBody)
	return response, err
}

// doPostEncodedForm performs a POST request with form-encoded data.
func (us *Uspacy) doPostEncodedForm(ctx context.Context, url string, values url.Values, opts ...RequestOption) ([]byte, error) {
	requestHeaders := mergeHeaders(opts, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	})

	response, _, err := us.doRaw(ctx, url, http.MethodPost, requestHeaders, []byte(values.Encode()))
	return response, err
}

// doDeleteEmptyHeaders performs a DELETE request with default headers and optional JSON body.
func (us *Uspacy) doDeleteEmptyHeaders(ctx context.Context, url string, body any, opts ...RequestOption) (int, error) {
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	requestHeaders := mergeHeaders(opts)
	_, code, err := us.doRaw(ctx, url, http.MethodDelete, requestHeaders, jsonBody)
	return code, err
}

// mergeHeaders builds a header map from default JSON headers plus any request options.
func mergeHeaders(opts []RequestOption, defaults ...map[string]string) map[string]string {
	requestHeaders := make(map[string]string)
	maps.Copy(requestHeaders, headersMap)
	for _, h := range defaults {
		maps.Copy(requestHeaders, h)
	}
	cfg := newRequestConfig(opts...)
	maps.Copy(requestHeaders, cfg.headers)
	return requestHeaders
}

// doPostFormData performs a multipart form POST request with files and text parameters.
// Returns error if no files are provided or if all provided files are invalid.
func (us *Uspacy) doPostFormData(ctx context.Context, url string, textParams map[string]string, files map[string]io.ReadCloser) ([]byte, error) {
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

	response, _, err := us.doRaw(ctx, url, http.MethodPost, headers, body.Bytes())
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
