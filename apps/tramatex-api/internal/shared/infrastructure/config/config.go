package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	DB       DatabaseConfig
	JWT      JWTConfig
	Security SecurityConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	AccessTokenTTL   string
	RefreshTokenTTL  string
	RefreshTokenTTL2 string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	JWTSecret         string
	JWTAccessTTL      string
	JWTRefreshTTL     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (not required in production)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "tramatex"),
			Password: getEnv("DB_PASSWORD", "tramatex"),
			Database: getEnv("DB_NAME", "tramatex_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:           getEnvRequired("JWT_SECRET"),
			AccessTokenTTL:   getEnv("JWT_ACCESS_TOKEN_TTL", "15m"),
			RefreshTokenTTL:  getEnv("JWT_REFRESH_TOKEN_TTL", "7d"),
			RefreshTokenTTL2: getEnv("JWT_REFRESH_TOKEN_TTL", "7d"),
		},
		Security: SecurityConfig{
			JWTSecret:     getEnvRequired("JWT_SECRET"),
			JWTAccessTTL:  getEnv("JWT_ACCESS_TOKEN_TTL", "15m"),
			JWTRefreshTTL: getEnv("JWT_REFRESH_TOKEN_TTL", "7d"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvRequired gets a required environment variable
func getEnvRequired(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

// DSN returns PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// Connect opens a database connection to PostgreSQL
// Returns a GORM DB instance ready for use
func (c *DatabaseConfig) Connect() (*gorm.DB, error) {
	dsn := c.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
