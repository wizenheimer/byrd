package config

import "time"

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
	ScreenshotServiceAPI string
	OpenAIKey            string
	ResendAPIKey         string
}

func Load() (*Config, error) {
	return nil, nil
}
