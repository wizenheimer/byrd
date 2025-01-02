package config

import "os"

// IsDevelopment returns true if the environment is set to development
func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
