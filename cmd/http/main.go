package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	handlers "github.com/solrac97gr/telegram-followers-checker/cmd/http/handlers"
	"github.com/solrac97gr/telegram-followers-checker/database"
	"github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
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

	// Initialize MongoDB client
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize components
	repo, err := database.NewMongoRepository(mongoClient)
	if err != nil {
		log.Fatalf("Error creating MongoDB repository: %v", err)
	}
	fm := filemanager.NewFileManager()
	telegramExtractor := telegram.NewTelegramExtractor()
	rutubeExtractor := rutube.NewRutubeExtractor()
	vkExtractor := vk.NewVKExtractor()
	instagramExtractor := instagram.NewInstagramExtractor()

	hdl := handlers.NewHandlers(
		repo,
		fm,
		telegramExtractor,
		rutubeExtractor,
		vkExtractor,
		instagramExtractor,
	)

	// Serve static files from the root directory
	fiberApp.Static("/", "./public")

	// Register handlers
	fiberApp.Post("/upload", hdl.UploadHandler)
	fiberApp.Get("/download", hdl.DownloadHandler)
	fiberApp.Post("/estimate-time", hdl.EstimateTimeHandler)

	log.Fatal(fiberApp.Listen(":3000"))
}
