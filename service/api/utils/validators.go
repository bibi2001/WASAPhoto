/*
The validatos.go contains the useful validation methods that are used the api subdirectory.
*/

package utils

import (
	"regexp"
)

func Validate(username string) bool {
	// Username can contain any numbers, letters and  also underscores
	// The length should be between 3 and 16 characters.
	pattern := `^[a-zA-Z0-9_][a-zA-Z0-9_]{2,16}$`

	// Compile the regular expression pattern.
	regex := regexp.MustCompile(pattern)

	// Use the regex to match the username.
	return regex.MatchString(username)
}
