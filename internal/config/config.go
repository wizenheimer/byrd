package config

import (
	"time"

	"github.com/wizenheimer/iris/pkg/utils/parser"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Services ServicesConfig
}

type ServerConfig struct {
	Port        string
	IdleTimeout time.Duration
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

	// Load server configuration.
	serverConfig := LoadServerConfig()

	// Load database configuration.
	databaseConfig := LoadDatabaseConfig()

	// Load storage configuration.
	storageConfig := LoadStorageConfig()

	// Load services configuration.
	servicesConfig := LoadServicesConfig()

	return &Config{
		Server:   serverConfig,
		Database: databaseConfig,
		Storage:  storageConfig,
		Services: servicesConfig,
	}, nil
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		// Port is set to the value of the SERVER_PORT environment variable, or "8080" if the variable is not set.
		Port: GetEnv("SERVER_PORT", "8080", parser.StrParser),

		// IdleTimeout is set to the value of the IDLE_TIMEOUT environment variable, or 300 seconds if the variable is not set.
		IdleTimeout: time.Duration(GetEnv("IDLE_TIMEOUT", 300, parser.IntParser)) * time.Second,
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
		ScreenshotServiceAPIKey: GetEnv("ScreenshotServiceAPIKey", "", parser.StrParser),
		OpenAIKey:               GetEnv("OpenAIKey", "", parser.StrParser),
		ResendAPIKey:            GetEnv("ResendAPIKey", "", parser.StrParser),
	}
}
