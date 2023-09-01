/*
The validatos.go contains the useful validation methods that are used the api subdirectory.
*/

package utils

import (
	"regexp"
	"strconv"
)

func ValidateCommentText(text string) bool {
	// Comment text can have maximum of 20 characters
	return len(text) < 21
}

func ValidateCaption(caption string) bool {
	// Caption can have a maximum of 90 characters
	return len(caption) < 91
}

func ValidateUsername(username string) bool {
	// Username can contain any numbers, letters and  also underscores
	// The length should be between 3 and 16 characters.
	pattern := `^[a-zA-Z0-9_][a-zA-Z0-9_]{2,16}$`

	// Compile the regular expression pattern.
	regex := regexp.MustCompile(pattern)

	// Use the regex to match the username.
	return regex.MatchString(username)
}

func ToInt64(s string) int64 {
	// Convert the user string to int64
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1

	}
	return res
}
