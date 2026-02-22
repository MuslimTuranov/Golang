package app

import (
	"Assignment2/pkg/modules"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() *modules.AppConfig {
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load("../.env")
	}

	ttlMin, _ := strconv.Atoi(getEnv("JWT_TTL_MIN", "60"))
	cfg := &modules.AppConfig{
		Port:      getEnv("APP_PORT", "8080"),
		APIKey:    getEnv("API_KEY", ""),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret"),
		JWTTTL:    time.Duration(ttlMin) * time.Minute,
		DB: modules.PostgreConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			Username:    getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", "postgres"),
			DBName:      getEnv("DB_NAME", "mydb"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			ExecTimeout: 5 * time.Second,
		},
	}

	if cfg.APIKey == "" {
		log.Println("Warning: API_KEY is empty")
	}
	return cfg
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
