package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	ServerPort string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_rewards"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "hdgi4u3ti4bot45t549t45t945bt945bt94t94t"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET не установлен")
	}

	if len(config.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET должен быть минимум 32 символа")
	}

	return config, nil
}

// GetDSN возвращает строку подключения к PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// GetPort возвращает порт сервера как int
func (c *Config) GetPort() int {
	port, err := strconv.Atoi(c.ServerPort)
	if err != nil {
		return 8080
	}
	return port
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
