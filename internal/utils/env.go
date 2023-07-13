package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	prod := os.Getenv("PROD")

	if prod != "true" {
		envPath := filepath.Join(".env")

		err := godotenv.Load(envPath)
		if err != nil {
			return err
		}

	}

	return nil
}
