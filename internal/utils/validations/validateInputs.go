package validations

import "errors"

func IsStringEmpty(str string) error {
	if str == "" {
		return errors.New("string is empty")
	}
	return nil
}

func IsPasswordValid(str string) error {
	if len(str) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}
	return nil
}
