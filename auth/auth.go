package auth

import (
	"errors"
	"net/http"
)

// Extracts the API key from the HTTP headers
func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("X-API-KEY")

	if apiKey == "" {
		return "", errors.New("no auth info found")
	}

	return apiKey, nil
}
