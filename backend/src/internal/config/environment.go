// ./src/internal/config/environment.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// List of possible .env files in order of precedence (later files override earlier ones)
	var envFiles []string

	// Check if ENV is set, if so use that to infer the environment
	switch os.Getenv("ENV") {
	case "development":
		envFiles = []string{
			".env",
			".env.local",
			".env.development",
		}
	case "test":
		envFiles = []string{
			".env",
			".env.test",
		}

	case "production":
		envFiles = []string{
			".env",
			".env.production",
		}
	default:
		envFiles = []string{
			".env",
			".env.local",
			".env.development",
			".env.test",
			".env.production",
		}
	}

	var loadedFiles []string
	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Overload(envFile); err != nil {
				return fmt.Errorf("error loading %s: %w", envFile, err)
			}
			loadedFiles = append(loadedFiles, envFile)
		}
	}

	if len(loadedFiles) == 0 {
		return fmt.Errorf("no environment files found for environment %s", os.Getenv("ENV"))
	}

	return nil
}

func GetEnv[T any](key string, fallback T, parser func(string) (T, error)) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	parsed, err := parser(value)
	if err != nil {
		return fallback
	}

	return parsed
}
