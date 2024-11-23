package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	AdminKey              string
	AdministratorKey      string
	JWTSecret             string
	BACKEND_URL           string
	RedisURL              string
	SegmindAPIKey         string
	SegmindSDAPI          string
	SegmindSamaritanAPI   string
	SegmindDreamshaperAPI string
	SegmindProtovisAPI    string
	AwsAccessKey          string
	AwsSecretKey          string
	AwsRegion             string
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
		Port:                  os.Getenv("PORT"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		AdminKey:              os.Getenv("ADMIN_KEY"),
		AdministratorKey:      os.Getenv("ADMINISTRATOR_KEY"),
		BACKEND_URL:           os.Getenv("BACKEND_URL"),
		RedisURL:              os.Getenv("REDIS_URL"),
		SegmindAPIKey:         os.Getenv("SEGMIND_API_KEY"),
		AwsAccessKey:          os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretKey:          os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AwsRegion:             os.Getenv("AWS_REGION"),
		SegmindSDAPI:          "https://api.segmind.com/v1/stable-diffusion-3.5-large-txt2img",
		SegmindProtovisAPI:    "https://api.segmind.com/v1/sdxl1.0-protovis-lightning",
		SegmindSamaritanAPI:   "https://api.segmind.com/v1/sdxl1.0-samaritan-3d",
		SegmindDreamshaperAPI: "https://api.segmind.com/v1/sdxl1.0-dreamshaper-lightning",
	}
}
