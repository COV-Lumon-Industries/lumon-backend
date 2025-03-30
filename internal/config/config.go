package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	Port          string
	SecretKey     []byte
	TokenDuration time.Duration
	GeminiAPIKey  string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBPort:        os.Getenv("DB_PORT"),
		Port:          os.Getenv("PORT"),
		SecretKey:     []byte(os.Getenv("SECRET")),
		TokenDuration: time.Duration(GetInt("TOKEN_EXPIRE_TIME", 24)) * time.Hour,
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

func GetInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("%s: %s", key, err)
			return fallback
		}
		return i
	}
	return fallback
}
