package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	Database DatabaseConfig
	AWS      AWSConfig
	App      AppConfig
	API      APIConfig
	Log      LogConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Webhook  WebhookConfig
}

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

// APIConfig holds API server configuration.
type APIConfig struct {
	Host    string
	Port    string
	BaseURL string
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	URL      string
}

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	Secret             string
	ExpiryHours        int
	RefreshExpiryHours int
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

// WebhookConfig holds webhook configuration.
type WebhookConfig struct {
	// DefaultSecret is used for generic providers or when provider-specific secret is not set
	DefaultSecret string
	// ProviderSecrets maps provider names to their specific secrets
	// Example: "opta" -> "opta-secret-key", "statsbomb" -> "statsbomb-secret-key"
	ProviderSecrets map[string]string
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string
	Format string
}

// AWSConfig holds AWS configuration.
type AWSConfig struct {
	Region           string
	AccessKeyID      string
	SecretAccessKey  string
	S3Bucket         string
	CloudFrontDomain string
}

const (
	// EnvironmentProduction represents the production environment.
	EnvironmentProduction = "production"
)

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error as it might not exist in production).
	//nolint:errcheck // .env file is optional, especially in production
	godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "footie"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("APP_ENV", "development"),
		},
		API: APIConfig{
			Host:    getEnv("API_HOST", "0.0.0.0"),
			Port:    getEnv("API_PORT", "8088"),
			BaseURL: getEnv("API_BASE_URL", "http://localhost:8088"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			Name:     getEnv("DATABASE_NAME", "footie"),
			User:     getEnv("DATABASE_USER", "footie_user"),
			Password: getEnv("DATABASE_PASSWORD", ""),
			SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key-change-this"),
			ExpiryHours:        getEnvAsInt("JWT_EXPIRY_HOURS", 24),
			RefreshExpiryHours: getEnvAsInt("JWT_REFRESH_EXPIRY_HOURS", 168),
		},
		CORS: CORSConfig{
			AllowedOrigins:   strings.Split(getEnv("CORS_ORIGINS", "http://localhost:4200"), ","),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		AWS: AWSConfig{
			Region:           getEnv("AWS_REGION", "eu-west-1"),
			AccessKeyID:      getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:         getEnv("AWS_S3_BUCKET", ""),
			CloudFrontDomain: getEnv("AWS_CLOUDFRONT_DOMAIN", ""),
		},
		Webhook: WebhookConfig{
			DefaultSecret: getEnv("WEBHOOK_SECRET", ""), // Default secret for generic providers
			ProviderSecrets: parseProviderSecrets(),      // Parse provider-specific secrets
		},
	}

	// Build DATABASE_URL if not provided
	if cfg.Database.URL == "" {
		cfg.Database.URL = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration values are set.
func (c *Config) Validate() error {
	if c.Database.Password == "" && c.App.Environment == EnvironmentProduction {
		return fmt.Errorf("DATABASE_PASSWORD is required in production")
	}

	if c.JWT.Secret == "your-secret-key-change-this" && c.App.Environment == EnvironmentProduction {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	return nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.App.Environment == EnvironmentProduction
}

// Helper functions.

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// parseProviderSecrets parses provider-specific webhook secrets from environment variables.
// Format: WEBHOOK_SECRET_<PROVIDER_NAME>=secret_value
// Example: WEBHOOK_SECRET_OPTA=opta-secret-key
//          WEBHOOK_SECRET_STATSBOMB=statsbomb-secret-key
func parseProviderSecrets() map[string]string {
	secrets := make(map[string]string)
	prefix := "WEBHOOK_SECRET_"

	// Iterate through all environment variables
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Check if it's a provider-specific secret
		if strings.HasPrefix(key, prefix) {
			// Extract provider name (e.g., "OPTA" from "WEBHOOK_SECRET_OPTA")
			providerName := strings.TrimPrefix(key, prefix)
			// Normalize to lowercase for consistency
			providerName = strings.ToLower(providerName)
			secrets[providerName] = value
		}
	}

	return secrets
}
