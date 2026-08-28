package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters for the application
type Config struct {
	AppEnv          string
	Port            string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	JWTExpiry       time.Duration
	HospitalAHISURL string
}

// Load loads configuration from environment variables or .env file
func Load() *Config {
	// Try loading .env if available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil || jwtExpiryHours <= 0 {
		jwtExpiryHours = 24
	}

	return &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "hospital_db"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecretjwtkeyforhospitalmiddleware123456"),
		JWTExpiry:       time.Duration(jwtExpiryHours) * time.Hour,
		HospitalAHISURL: getEnv("HOSPITAL_A_HIS_URL", "https://hospital-a.api.co.th"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
