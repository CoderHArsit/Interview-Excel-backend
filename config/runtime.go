package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Runtime struct {
	AppEnv             string
	Port               string
	DatabaseURL        string
	CookieDomain       string
	CookieSecure       bool
	CorsAllowedOrigins []string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	RedisEnabled       bool
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	RedisUseTLS        bool
	RazorpayKey        string
	RazorpaySecret     string
}

var (
	runtimeConfig Runtime
	runtimeOnce   sync.Once
)

func RuntimeConfig() Runtime {
	runtimeOnce.Do(func() {
		_ = godotenv.Load()

		appEnv := getEnv("APP_ENV", "development")
		port := getEnv("PORT", "8080")

		runtimeConfig = Runtime{
			AppEnv:             appEnv,
			Port:               port,
			DatabaseURL:        strings.TrimSpace(os.Getenv("DATABASE_URL")),
			CookieDomain:       strings.TrimSpace(os.Getenv("COOKIE_DOMAIN")),
			CookieSecure:       getEnvBool("COOKIE_SECURE", appEnv == "production"),
			CorsAllowedOrigins: getEnvCSV("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3010"}),
			GoogleClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
			GoogleClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", fmt.Sprintf("http://localhost:%s/google_callback", port)),
			RedisEnabled:       resolveRedisEnabled(),
			RedisAddr:          strings.TrimSpace(os.Getenv("REDIS_ADDR")),
			RedisPassword:      os.Getenv("REDIS_PASSWORD"),
			RedisDB:            getEnvInt("REDIS_DB", 0),
			RedisUseTLS:        getEnvBool("REDIS_USE_TLS", false),
			RazorpayKey:        strings.TrimSpace(os.Getenv("RAZORPAY_KEY")),
			RazorpaySecret:     strings.TrimSpace(os.Getenv("RAZORPAY_SECRET")),
		}
	})

	return runtimeConfig
}

func DatabaseDSN() string {
	if databaseURL := RuntimeConfig().DatabaseURL; databaseURL != "" {
		return databaseURL
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""),
		getEnv("DB_USER", ""),
		getEnv("DB_NAME", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_SSLMODE", "disable"),
	)
}

func RedisTLSConfig() *tls.Config {
	if !RuntimeConfig().RedisUseTLS {
		return nil
	}

	return &tls.Config{MinVersion: tls.VersionTLS12}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getEnvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvCSV(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	if len(origins) == 0 {
		return fallback
	}

	return origins
}

func resolveRedisEnabled() bool {
	value, ok := os.LookupEnv("REDIS_ENABLED")
	if !ok {
		return strings.TrimSpace(os.Getenv("REDIS_ADDR")) != ""
	}

	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return false
	}

	return parsed
}
