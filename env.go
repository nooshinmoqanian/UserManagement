package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Port          string
	JWTSecret     string
	PGURL         string
	RedisAddr     string
	RedisDB       int
	OTPExpiresIn  time.Duration
	OTPMaxReq     int
	OTPWindow     time.Duration
	EnableSwagger bool
}

// LoadConfig reads configuration from environment variables
func LoadConfig() *Config {
	c := &Config{
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-insecure-secret"),
		PGURL:         getEnv("PG_URL", "postgres://otpuser:otppass@localhost:5432/otpdb?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		OTPExpiresIn:  getEnvDuration("OTP_EXPIRES_IN", 2*time.Minute),
		OTPMaxReq:     getEnvInt("OTP_MAX_REQUESTS", 3),
		OTPWindow:     getEnvDuration("OTP_WINDOW", 10*time.Minute),
		EnableSwagger: getEnv("ENABLE_SWAGGER", "true") == "true",
	}

	if c.JWTSecret == "dev-insecure-secret" {
		log.Println("WARN: Using insecure default JWT secret. Set JWT_SECRET in production!")
	}

	return c
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var i int
		_, err := fmt.Sscanf(val, "%d", &i)
		if err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		d, err := time.ParseDuration(val)
		if err == nil {
			return d
		}
	}
	return defaultVal
}
