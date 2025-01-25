package config

import "os"

// IsDevelopment returns true if the environment is set to development
// or if the ENV_PROFILE is set to development
// or if the ENV_LOG_LEVEL is set to debug
func IsDevelopment() bool {
	return os.Getenv("ENV") == "development" || os.Getenv("ENV_PROFILE") == "development" || os.Getenv("ENV_LOG_LEVEL") == "debug"
}
