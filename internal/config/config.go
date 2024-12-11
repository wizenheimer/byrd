package config

import (
	"time"

	"github.com/wizenheimer/iris/pkg/utils/parser"
)

type Config struct {
	Environment EnvironmentConfig
	Server      ServerConfig
	Database    DatabaseConfig
	Storage     StorageConfig
	Services    ServicesConfig
}

type ServerConfig struct {
	Port                string
	IdleTimeout         time.Duration
	ShutdownTimeout     time.Duration
	ShutdownMaxAttempts int
}

type EnvironmentConfig struct {
	EnvProfile  string
	EnvLogLevel string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type StorageConfig struct {
	Type      string // "local" or "r2"
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

type ServicesConfig struct {
	ScreenshotServiceAPIKey string
	OpenAIKey               string
	ResendAPIKey            string
}

func Load() (*Config, error) {
	// Load environment variables from a .env file.
	if err := LoadEnv(); err != nil {
		return nil, err
	}

	// Load environment configuration.
	environmentConfig := LoadEnvironmentConfig()

	// Load server configuration.
	serverConfig := LoadServerConfig()

	// Load database configuration.
	databaseConfig := LoadDatabaseConfig()

	// Load storage configuration.
	storageConfig := LoadStorageConfig()

	// Load services configuration.
	servicesConfig := LoadServicesConfig()

	return &Config{
		Environment: environmentConfig,
		Server:      serverConfig,
		Database:    databaseConfig,
		Storage:     storageConfig,
		Services:    servicesConfig,
	}, nil
}

func LoadEnvironmentConfig() EnvironmentConfig {
	return EnvironmentConfig{
		// EnvProfile is set to the value of the ENV_PROFILE environment variable, or "development" if the variable is not set.
		EnvProfile: GetEnv("ENV_PROFILE", "development", parser.StrParser),
		// EnvLogLevel is set to the value of the ENV_LOG_LEVEL environment variable, or "info" if the variable is not set.
		EnvLogLevel: GetEnv("ENV_LOG_LEVEL", "info", parser.StrParser),
	}
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		// Port is set to the value of the SERVER_PORT environment variable, or "8080" if the variable is not set.
		Port: GetEnv("SERVER_PORT", "8080", parser.StrParser),

		// IdleTimeout is set to the value of the IDLE_TIMEOUT environment variable, or 300 seconds if the variable is not set.
		IdleTimeout: time.Duration(GetEnv("SERVER_IDLE_TIMEOUT", 300, parser.IntParser)) * time.Second,

		// ShutdownTimeout is set to the value of the SHUTDOWN_TIMEOUT environment variable, or 10 seconds if the variable is not set.
		ShutdownTimeout: time.Duration(GetEnv("SERVER_SHUTDOWN_TIMEOUT", 10, parser.IntParser)) * time.Second,

		// ShutdownMaxAttempts is set to the value of the SHUTDOWN_MAX_ATTEMPTS environment variable, or 3 if the variable is not set.
		ShutdownMaxAttempts: GetEnv("SERVER_SHUTDOWN_MAX_ATTEMPTS", 3, parser.IntParser),
	}
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		// Driver is set to the value of the DB_DRIVER environment variable, or "postgres" if the variable is not set.
		Driver: GetEnv("DB_DRIVER", "postgres", parser.StrParser),
		// Host is set to the value of the DB_HOST environment variable, or "localhost" if the variable is not set.
		Host: GetEnv("DB_HOST", "localhost", parser.StrParser),
		// Port is set to the value of the DB_PORT environment variable, or "5432" if the variable is not set.
		Port: GetEnv("DB_PORT", "5432", parser.StrParser),
		// User is set to the value of the DB_USER environment variable, or "postgres" if the variable is not set.
		User: GetEnv("DB_USER", "postgres", parser.StrParser),
		// Password is set to the value of the DB_PASSWORD environment variable, or "postgres" if the variable is not set.
		Password: GetEnv("DB_PASSWORD", "postgres", parser.StrParser),
		// Database is set to the value of the DB_DATABASE environment variable, or "postgres" if the variable is not set.
		Database: GetEnv("DB_DATABASE", "postgres", parser.StrParser),
	}
}

func LoadStorageConfig() StorageConfig {
	return StorageConfig{
		// Type is set to the value of the STORAGE_TYPE environment variable, or "local" if the variable is not set.
		Type: GetEnv("STORAGE_TYPE", "local", parser.StrParser),
		// AccessKey is set to the value of the STORAGE_ACCESS_KEY environment variable, or "" if the variable is not set.
		AccessKey: GetEnv("STORAGE_ACCESS_KEY", "", parser.StrParser),
		// SecretKey is set to the value of the STORAGE_SECRET_KEY environment variable, or "" if the variable is not set.
		SecretKey: GetEnv("STORAGE_SECRET_KEY", "", parser.StrParser),
		// Bucket is set to the value of the STORAGE_BUCKET environment variable, or "" if the variable is not set.
		Bucket: GetEnv("STORAGE_BUCKET", "", parser.StrParser),
		// Region is set to the value of the STORAGE_REGION environment variable, or "" if the variable is not set.
		Region: GetEnv("STORAGE_REGION", "", parser.StrParser),
	}
}

func LoadServicesConfig() ServicesConfig {
	return ServicesConfig{
		ScreenshotServiceAPIKey: GetEnv("SCREENSHOT_API_KEY", "", parser.StrParser),
		OpenAIKey:               GetEnv("OPENAI_API_KEY", "", parser.StrParser),
		ResendAPIKey:            GetEnv("RESEND_API_KEY", "", parser.StrParser),
	}
}
