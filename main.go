package main

import (
	"karmaclips/bedrock"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// func main() {
// 	err := godotenv.Load()
// 	opts := &slog.HandlerOptions{
// 		AddSource: true,
// 		Level:     slog.LevelDebug,
// 	}
// 	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
// 	if err != nil {
// 		logger.Error("unable to load .env")
// 	}

// 	app := routes.Routes()
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "*",
// 	}))
// 	app.Listen(":6969")
// }

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

	// bedrock.PromptModel("\n\nHuman: explain black holes to 8th graders\n\nAssistant:", 0.1, 0.9, 50)
	bedrock.PromptModelStream("\n\nHuman: explain black holes to 8th graders\n\nAssistant:", 0.1, 0.9, 50)
}
