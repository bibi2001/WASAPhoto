package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(r *http.Request) (int64, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if the Authorization header is empty or does not start with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New("invalid authorization header")
	}

	// Extract the user identifier (token) from the header
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return ToInt64(token), nil
}
