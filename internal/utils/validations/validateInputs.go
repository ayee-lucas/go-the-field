package validations

import (
	"errors"
	"regexp"
)

// IsStringEmpty checks if a string is empty.
func IsStringEmpty(str string) error {
	if str == "" {
		return errors.New("string is empty")
	}
	return nil
}

// IsPasswordValid checks if a password meets the complexity requirements.
func IsPasswordValid(str string) error {
	if len(str) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*()]`).MatchString

	if !hasUppercase(str) {
		return errors.New("Password must contain at least one uppercase letter")
	}

	if !hasLowercase(str) {
		return errors.New("Password must contain at least one lowercase letter")
	}

	if !hasNumber(str) {
		return errors.New("Password must contain at least one number")
	}

	if !hasSpecialChar(str) {
		return errors.New("Password must contain at least one special character")
	}

	return nil
}

// IsEmailValid checks if an email address is valid.
func IsEmailValid(str string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(str) {
		return errors.New("Email is not valid")
	}

	return nil
}

// IsUsernameValid checks if a username is valid.
func IsUsernameValid(str string) error {
	// Define the regular expression pattern for a valid username
	usernameRegex := `^[a-zA-Z0-9._]{3,20}$`

	// Compile the regular expression pattern
	re := regexp.MustCompile(usernameRegex)

	// Check if the username matches the regular expression pattern
	if !re.MatchString(str) {
		return errors.New("Username is not valid")
	}

	return nil
}
