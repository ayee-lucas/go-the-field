package validations

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(password, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err
}
