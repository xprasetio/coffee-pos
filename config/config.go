package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	// App
	AppPort string
	AppEnv  string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Midtrans
	MidtransServerKey string
	MidtransClientKey string
	MidtransEnv       string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (ignore error for production)
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		AppEnv:            getEnv("APP_ENV", "development"),
		DBHost:            getEnv("DB_HOST", ""),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBName:            getEnv("DB_NAME", ""),
		DBUser:            getEnv("DB_USER", ""),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		JWTExpiryHours:    24,
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransEnv:       getEnv("MIDTRANS_ENV", "sandbox"),
	}

	// Parse JWTExpiryHours
	if expiryHours := getEnv("JWT_EXPIRY_HOURS", ""); expiryHours != "" {
		if hours, err := strconv.Atoi(expiryHours); err == nil {
			cfg.JWTExpiryHours = hours
		}
	}

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks required configuration fields
func (c *Config) validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// MysqlDSN returns MySQL connection string
func (c *Config) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// RedisAddr returns Redis address in host:port format
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
