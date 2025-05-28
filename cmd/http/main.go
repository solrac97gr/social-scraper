package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	handlers "github.com/solrac97gr/telegram-followers-checker/cmd/http/handlers"
)

const (
	EchoVar = "ECHO_VAR"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Print("Environment variables loaded successfully echo var:", os.Getenv(EchoVar))
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	// Serve static files from the root directory
	fiberApp.Static("/", "./public")

	// Register handlers
	fiberApp.Post("/upload", handlers.UploadHandler)
	fiberApp.Get("/download", handlers.DownloadHandler)
	fiberApp.Post("/estimate-time", handlers.EstimateTimeHandler)

	log.Fatal(fiberApp.Listen(":3000"))
}
