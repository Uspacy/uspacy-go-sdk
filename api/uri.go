package api

// jsonHeaders returns a new map with the default headers of a JSON request. Each call
// returns a fresh map, so callers can add to it without affecting other requests.
func jsonHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}
