// ./src/internal/config/config.go
package config

import (
	"fmt"
	"time"

	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type Config struct {
	ServiceName string
	Environment EnvironmentConfig
	Server      ServerConfig
	Database    DatabaseConfig
	Storage     StorageConfig
	Services    ServicesConfig
	Workflow    WorkflowConfig
}

type ServerConfig struct {
	// Port is the port on which the server listens for incoming requests.
	Port string

	// Limits the number of requests for all endpoints.
	GlobalRequestsPerMinute int

	// Limits the number of requests for workspace creation and deletion.
	// Defaults to 1 request per minute.
	WorkspaceCDRequestsPerMinute int

	// Limits the number of requests for competitor creation and deletion.
	// Defaults to 5 requests per second.
	CompetitorCDRequestsPerSecond int

	// Limits the number of requests for page creation and deletion.
	// Defaults to 10 requests per second.
	PageCDRequestsPerSecond int

	// UserCDRequestsPerSecond limits the number of requests for user creation and deletion.
	// Defaults to 1 request per second.
	UserCDRequestsPerSecond int

	// Limit the number of report created or dispatched per minute.
	// Defaults to 10 requests per minute.
	ReportCDPerMinute int

	// CorsAllowedOrigins is a comma-separated list of origins that are allowed to make requests to the server.
	CorsAllowedOrigins string

	// IdleTimeout is the time after which the server waits for the next request before shutting down.
	IdleTimeout time.Duration

	// ShutdownTimeout is the maximum amount of time to wait for the server to shut down gracefully.
	ShutdownTimeout time.Duration

	// ShutdownMaxAttempts is the maximum number of attempts to shut down the server gracefully.
	ShutdownMaxAttempts int
}

type EnvironmentConfig struct {
	EnvProfile  string
	EnvLogLevel string
}

type DatabaseConfig struct {
	// Driver for the database
	Driver string

	// Host and port of the database
	Host string
	Port string

	// Credentials for the database
	User     string
	Password string

	// Name of the database
	Database string

	// Connection string for the database
	ConnectionString string
}

type StorageConfig struct {
	BaseEndpoint string // s3 base endpoint (for r2 it's infered from the account id)
	Type         string // "local" or "r2"
	AccessKey    string
	SecretKey    string
	Bucket       string
	Region       string
	AccountId    string
}

type ServicesConfig struct {
	ClerkAPIKey                  string
	HightlightProjectID          string
	ScreenshotServiceAPIKey      string
	ScreenshotServiceOrigin      string
	ScreenshotServiceQPS         float64
	OpenAIKey                    string
	ResendAPIKey                 string
	ResendNotificationEmail      string
	PostHogAPIKey                string
	ManagementAPIKey             string
	ManagementAPIRefreshInterval time.Duration
}

type WorkflowConfig struct {
	// RedisURL is the URL of the Redis server
	RedisURL string
	// RedisAddr is the address of the Redis server
	RedisAddr string
	// RedisPassword is the password for the Redis server
	RedisPassword string
	// RedisDB is the database number for the Redis server
	RedisDB int
	// WorkflowTTL is the time-to-live for the workflow
	WorkflowTTL time.Duration
	// ExecutorParallelism is the number of parallel jobs
	ReportExecutorParallelism int
	// ScreenshotExecutorParallelism is the number of parallel jobs
	ScreenshotExecutorParallelism int
	// ReportExecutorLowerBound is the lower bound for the executor
	ReportExecutorLowerBound int
	// ExecutorLowerBound is the lower bound for the executor
	ScreenshotExecutorLowerBound int
	// ReportExecutorUpperBound is the upper bound for the executor
	ReportExecutorUpperBound int
	// ExecutorUpperBound is the upper bound for the executor
	ScreenshotExecutorUpperBound int
}

func Load() (*Config, error) {
	// Load environment variables from a .env file
	envProfile := GetEnv("ENV_PROFILE", "development", utils.StrParser)
	// LoadEnv is a function that loads environment variables from a .env file
	// It returns an error only if both
	// - The file is not found
	// - The environment profile is not "production"
	if err := LoadEnv(); err != nil && envProfile != "production" {
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

	// Load workflow configuration.
	workflowConfig := LoadWorkflowConfig()

	return &Config{
		ServiceName: "byrd-go",
		Environment: environmentConfig,
		Server:      serverConfig,
		Database:    databaseConfig,
		Storage:     storageConfig,
		Services:    servicesConfig,
		Workflow:    workflowConfig,
	}, nil
}

func LoadEnvironmentConfig() EnvironmentConfig {
	return EnvironmentConfig{
		// EnvProfile is set to the value of the ENV_PROFILE environment variable, or "development" if the variable is not set.
		EnvProfile: GetEnv("ENV_PROFILE", "development", utils.StrParser),
		// EnvLogLevel is set to the value of the ENV_LOG_LEVEL environment variable, or "info" if the variable is not set.
		EnvLogLevel: GetEnv("ENV_LOG_LEVEL", "info", utils.StrParser),
	}
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		// Port is set to the value of the SERVER_PORT environment variable, or "8080" if the variable is not set.
		Port: fmt.Sprintf(":%s", GetEnv("PORT", "10000", utils.StrParser)),

		// RequestsPerMinute is set to the value of the REQUESTS_PER_MINUTE environment variable, or 100 if the variable is not set.
		GlobalRequestsPerMinute: GetEnv("GLOBAL_REQUESTS_PER_MINUTE", 100, utils.IntParser),

		// WorkspaceCDRequestsPerSecond is set to the value of the WORKSPACE_CD_REQUESTS_PER_SECOND environment variable, or 1 if the variable is not set.
		WorkspaceCDRequestsPerMinute: GetEnv("WORKSPACE_CD_REQUESTS_PER_MINUTE", 1, utils.IntParser),

		// CompetitorCDRequestsPerSecond is set to the value of the COMPETITOR_CD_REQUESTS_PER_SECOND environment variable, or 5 if the variable is not set.
		CompetitorCDRequestsPerSecond: GetEnv("COMPETITOR_CD_REQUESTS_PER_SECOND", 5, utils.IntParser),

		// PageCDRequestsPerSecond is set to the value of the PAGE_CD_REQUESTS_PER_SECOND environment variable, or 10 if the variable is not set.
		PageCDRequestsPerSecond: GetEnv("PAGE_CD_REQUESTS_PER_SECOND", 10, utils.IntParser),

		// UserCDRequestsPerSecond is set to the value of the USER_CD_REQUESTS_PER_SECOND environment variable, or 1 if the variable is not set.
		UserCDRequestsPerSecond: GetEnv("USER_CD_REQUESTS_PER_SECOND", 1, utils.IntParser),

		// ReportCDPerMinute is set to the value of the REPORT_CD_PER_MINUTE environment variable, or 10 if the variable is not set.
		ReportCDPerMinute: GetEnv("REPORT_CD_PER_MINUTE", 10, utils.IntParser),

		// CorsAllowedOrigins is set to the value of the CORS_ALLOWED_ORIGINS environment variable, or "http://localhost:5173" if the variable is not set.
		CorsAllowedOrigins: GetEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173", utils.StrParser),

		// IdleTimeout is set to the value of the IDLE_TIMEOUT environment variable, or 300 seconds if the variable is not set.
		IdleTimeout: time.Duration(GetEnv("SERVER_IDLE_TIMEOUT", 300, utils.IntParser)) * time.Second,

		// ShutdownTimeout is set to the value of the SHUTDOWN_TIMEOUT environment variable, or 10 seconds if the variable is not set.
		ShutdownTimeout: time.Duration(GetEnv("SERVER_SHUTDOWN_TIMEOUT", 10, utils.IntParser)) * time.Second,

		// ShutdownMaxAttempts is set to the value of the SHUTDOWN_MAX_ATTEMPTS environment variable, or 3 if the variable is not set.
		ShutdownMaxAttempts: GetEnv("SERVER_SHUTDOWN_MAX_ATTEMPTS", 3, utils.IntParser),
	}
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		// Driver is set to the value of the DB_DRIVER environment variable, or "postgres" if the variable is not set.
		Driver: GetEnv("DB_DRIVER", "postgres", utils.StrParser),
		// Host is set to the value of the DB_HOST environment variable, or "localhost" if the variable is not set.
		Host: GetEnv("DB_HOST", "localhost", utils.StrParser),
		// Port is set to the value of the DB_PORT environment variable, or "5432" if the variable is not set.
		Port: GetEnv("DB_PORT", "5432", utils.StrParser),
		// User is set to the value of the DB_USER environment variable, or "postgres" if the variable is not set.
		User: GetEnv("DB_USER", "postgres", utils.StrParser),
		// Password is set to the value of the DB_PASSWORD environment variable, or "postgres" if the variable is not set.
		Password: GetEnv("DB_PASSWORD", "postgres", utils.StrParser),
		// Database is set to the value of the DB_DATABASE environment variable, or "postgres" if the variable is not set.
		Database: GetEnv("DB_DATABASE", "postgres", utils.StrParser),
		// ConnectionString is set to the value of the DB_CONNECTION_STRING environment variable, or "" if the variable is not set.
		ConnectionString: GetEnv("DB_CONNECTION_STRING", "", utils.StrParser),
	}
}

func LoadStorageConfig() StorageConfig {
	return StorageConfig{
		// BaseEndpoint is set to the value of the STORAGE_BASE_ENDPOINT environment variable, or "" if the variable is not set.
		BaseEndpoint: GetEnv("STORAGE_BASE_ENDPOINT", "", utils.StrParser),
		// Type is set to the value of the STORAGE_TYPE environment variable, or "local" if the variable is not set.
		Type: GetEnv("STORAGE_TYPE", "local", utils.StrParser),
		// AccessKey is set to the value of the STORAGE_ACCESS_KEY environment variable, or "" if the variable is not set.
		AccessKey: GetEnv("STORAGE_ACCESS_KEY", "", utils.StrParser),
		// SecretKey is set to the value of the STORAGE_SECRET_KEY environment variable, or "" if the variable is not set.
		SecretKey: GetEnv("STORAGE_SECRET_KEY", "", utils.StrParser),
		// Bucket is set to the value of the STORAGE_BUCKET environment variable, or "" if the variable is not set.
		Bucket: GetEnv("STORAGE_BUCKET", "", utils.StrParser),
		// Region is set to the value of the STORAGE_REGION environment variable, or "" if the variable is not set.
		Region: GetEnv("STORAGE_REGION", "", utils.StrParser),
		// AccountId is set to the value of the STORAGE_ACCOUNT_ID environment variable, or "" if the variable is not set.
		AccountId: GetEnv("STORAGE_ACCOUNT_ID", "", utils.StrParser),
	}
}

func LoadServicesConfig() ServicesConfig {
	return ServicesConfig{
		// ClerkAPIKey is set to the value of the CLERK_API_KEY environment variable, or "" if the variable is not set.
		ClerkAPIKey: GetEnv("CLERK_API_KEY", "", utils.StrParser),
		// HightlightProjectID is set to the value of the HIGHLIGHT_PROJECT_ID environment variable, or "" if the variable is not set.
		HightlightProjectID: GetEnv("HIGHLIGHT_PROJECT_ID", "", utils.StrParser),
		// ScreenshotServiceAPIKey is set to the value of the SCREENSHOT_API_KEY environment variable, or "" if the variable is not set.
		ScreenshotServiceAPIKey: GetEnv("SCREENSHOT_API_KEY", "", utils.StrParser),
		// ScreenshotServiceOrigin is set to the value of the SCREENSHOT_API_ORIGIN environment variable, or "" if the variable is not set.
		ScreenshotServiceOrigin: GetEnv("SCREENSHOT_API_ORIGIN", "", utils.StrParser),
		// ScreenshotServiceQPS is set to the value of the SCREENSHOT_API_QPS environment variable, or 0.667 if the variable is not set.
		ScreenshotServiceQPS: GetEnv("SCREENSHOT_API_QPS", 0.667, utils.Float64Parser),
		// OpenAIKey is set to the value of the OPENAI_API_KEY environment variable, or "" if the variable is not set.
		OpenAIKey: GetEnv("OPENAI_API_KEY", "", utils.StrParser),
		// ResendAPIKey is set to the value of the RESEND_API_KEY environment variable, or "" if the variable is not set.
		ResendAPIKey: GetEnv("RESEND_API_KEY", "", utils.StrParser),
		// ResendNotificationEmail is set to the value of the RESEND_NOTIFICATION_EMAIL environment variable, or "" if the variable is not set.
		ResendNotificationEmail: GetEnv("RESEND_NOTIFICATION_EMAIL", "hey@byrdhq.com", utils.StrParser),
		// PostHogAPIKey is set to the value of the POSTHOG_API_KEY environment variable, or "" if the variable is not set.
		PostHogAPIKey: GetEnv("POSTHOG_API_KEY", "", utils.StrParser),
		// ManagementAPIKey is set to the value of the MANAGEMENT_API_KEY environment variable, or "" if the variable is not set.
		ManagementAPIKey: GetEnv("MANAGEMENT_API_KEY", "", utils.StrParser),
		// ManagementAPIRefreshInterval is set to the value of the MANAGEMENT_API_REFRESH_INTERVAL environment variable, or 5 minutes if the variable is not set.
		ManagementAPIRefreshInterval: time.Duration(GetEnv("MANAGEMENT_API_REFRESH_INTERVAL", 5, utils.IntParser)) * time.Minute,
	}
}

func LoadWorkflowConfig() WorkflowConfig {
	// Load workflow configuration.
	return WorkflowConfig{
		// RedisURL is set to the value of the REDIS_URL environment variable, or "" if the variable is not set.
		RedisURL: GetEnv("REDIS_URL", "", utils.StrParser),
		// RedisAddr is set to the value of the REDIS_ADDR environment variable, or "" if the variable is not set.
		RedisAddr: GetEnv("REDIS_ADDR", "", utils.StrParser),
		// RedisPassword is set to the value of the REDIS_PASSWORD environment variable, or "" if the variable is not set.
		RedisPassword: GetEnv("REDIS_PASSWORD", "", utils.StrParser),
		// RedisDB is set to the value of the REDIS_DB environment variable, or 0 if the variable is not set.
		RedisDB: GetEnv("REDIS_DB", 0, utils.IntParser),
		// RedisConnectionStr is set to the value of the REDIS_CONNECTION_STR environment variable, or "" if the variable is not set.
		// WorkflowTTL is set to the value of the WORKFLOW_TTL environment variable, or 4 days if the variable is not set.
		WorkflowTTL: time.Duration(GetEnv("WORKFLOW_TTL", 4*24*60, utils.IntParser)) * time.Second,
		// ExecutorParallelism is set to the value of the EXECUTOR_PARALLELISM environment variable, or 10 if the variable is not set.
		ScreenshotExecutorParallelism: GetEnv("SCREENSHOT_EXECUTOR_PARALLELISM", 10, utils.IntParser),
		// ExecutorLowerBound is set to the value of the EXECUTOR_LOWER_BOUND environment variable, or 10 seconds if the variable is not set.
		ScreenshotExecutorLowerBound: GetEnv("SCREENSHOT_EXECUTOR_LOWER_BOUND", 10, utils.IntParser),
		// ExecutorUpperBound is set to the value of the EXECUTOR_UPPER_BOUND environment variable, or 20 seconds if the variable is not set.
		ScreenshotExecutorUpperBound: GetEnv("SCREENSHOT_EXECUTOR_UPPER_BOUND", 120, utils.IntParser),
		// ReportExecutorParallelism is set to the value of the REPORT_EXECUTOR_PARALLELISM environment variable, or 10 if the variable is not set.
		ReportExecutorParallelism: GetEnv("REPORT_EXECUTOR_PARALLELISM", 10, utils.IntParser),
		// ReportExecutorLowerBound is set to the value of the REPORT_EXECUTOR_LOWER_BOUND environment variable, or 10 seconds if the variable is not set.
		ReportExecutorLowerBound: GetEnv("REPORT_EXECUTOR_LOWER_BOUND", 10, utils.IntParser),
		// ReportExecutorUpperBound is set to the value of the REPORT_EXECUTOR_UPPER_BOUND environment variable, or 20 seconds if the variable is not set.
		ReportExecutorUpperBound: GetEnv("REPORT_EXECUTOR_UPPER_BOUND", 120, utils.IntParser),
	}
}
