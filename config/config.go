package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	AdminKey         string
	AdministratorKey string
	JWTSecret        string
	BACKEND_URL      string
	RedisURL         string
	RedisToken       string
	SegmindAPIKey    string
	SegmindSDAPI     string
	AwsAccessKey     string
	AwsSecretKey     string
	AwsRegion        string
}

func NewConfig() *Config {
	err := godotenv.Load()
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	if err != nil {
		logger.Error("unable to load .env")
	}
	return &Config{
		Port:             os.Getenv("PORT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		AdminKey:         os.Getenv("ADMIN_KEY"),
		AdministratorKey: os.Getenv("ADMINISTRATOR_KEY"),
		BACKEND_URL:      os.Getenv("BACKEND_URL"),
		RedisURL:         os.Getenv("REDIS_URL"),
		RedisToken:       os.Getenv("REDIS_TOKEN"),
		SegmindAPIKey:    os.Getenv("SEGMIND_API_KEY"),
		SegmindSDAPI:     os.Getenv("SEGMIND_SD_API"),
		AwsAccessKey:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretKey:     os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AwsRegion:        os.Getenv("AWS_REGION"),
	}
}
