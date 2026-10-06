package api

// headersMap holds the default headers of a JSON request. Helpers copy it before adding
// request options, so callers never mutate it.
var headersMap = map[string]string{
	"Content-Type": "application/json",
}
