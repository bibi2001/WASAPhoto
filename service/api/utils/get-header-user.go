package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(r *http.Request) (string, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if the Authorization header is empty or does not start with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Invalid Authorization header")
	}

	// Extract the user identifier (token) from the header
	userIdentifier := strings.TrimPrefix(authHeader, "Bearer ")

	return userIdentifier, nil
}
