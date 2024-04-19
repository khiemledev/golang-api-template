package util

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

// IsValidEmail checks if the email provided is valid by regex.
func IsValidEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// IsValidPassword checks if the password is valid. The valid password must be
// at least 8 characters long and contain at least 1 uppercase letter, 1 lowercase
// letter, and 1 number
func IsValidPassword(password string) bool {
	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	hasNumber, hasUpper, hasLower, hasSpecial := false, false, false, false
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasNumber && hasUpper && hasLower && hasSpecial
}
