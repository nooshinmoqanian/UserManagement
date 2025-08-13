package config

import (
	"log"
	"os"
	"time"
)

const (
	DefaultOTPExpiry = 2 * time.Minute
	DefaultOTPLimit  = 3
	DefaultOTPWindow = 10 * time.Minute
	DefaultJWTSecret = "dev-insecure-secret"
)

type Config struct {
	JWTSecret string
	Port      string
	PGURL     string
	RedisAddr string
	OTPExpire time.Duration
	OTPMaxReq int
	OTPWindow time.Duration
}

func Load() Config {
	return Config{
		JWTSecret: getEnv("JWT_SECRET", DefaultJWTSecret),
		Port:      getEnv("PORT", "8080"),
		PGURL:     mustEnv("PG_URL"),
		RedisAddr: mustEnv("REDIS_ADDR"),
		OTPExpire: DefaultOTPExpiry,
		OTPMaxReq: DefaultOTPLimit,
		OTPWindow: DefaultOTPWindow,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Printf("WARN: %s not set; using default", key)
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("ERROR: %s is required but not set", key)
	}
	return v
}
