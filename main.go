package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"karmaclips/routes"
)

func main() {
	err := godotenv.Load()
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	if err != nil {
		logger.Error("unable to load .env")
	}

	app := routes.Routes()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Listen(":6969")
}
