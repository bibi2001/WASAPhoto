package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(r *http.Request) (int64, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Log the received Authorization header
	fmt.Println("Received Authorization Header:", authHeader)

	// Check if the Authorization header is empty or does not start with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New("invalid or missing Authorization header")
	}

	// Extract the user identifier (token) from the header
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return ToInt64(token), nil
}
