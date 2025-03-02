package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	handlers "github.com/solrac97gr/telegram-followers-checker/cmd/http/handlers"
)

func main() {
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
