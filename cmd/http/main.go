package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/solrac97gr/telegram-followers-checker/app"
	handlers "github.com/solrac97gr/telegram-followers-checker/cmd/http/handlers"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/config"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/repository"
)

func main() {
	// Initialize MongoDB client
	mongoURI := "mongodb://localhost:27017"
	client, err := app.InitMongoClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB client: %v", err)
	}

	repo := repository.NewMongoRepository(client, "tgstats_db")

	// Load TGStat configuration
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load TGStat configuration: %v", err)
	}

	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	// Serve static files from the root directory
	fiberApp.Static("/", "./public")

	// Register handlers
	fiberApp.Post("/upload", handlers.UploadHandler(repo, config))
	fiberApp.Get("/download", handlers.DownloadHandler)
	fiberApp.Post("/estimate-time", handlers.EstimateTimeHandler)

	log.Fatal(fiberApp.Listen(":3000"))
}
