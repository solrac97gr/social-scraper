package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/app"
	"github.com/solrac97gr/telegram-followers-checker/config"
	"github.com/solrac97gr/telegram-followers-checker/database"
	instagram "github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	vk "github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	startAt := time.Now()
	// Check if input argument is provided
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the Excel file as an argument")
	}

	inputFile := os.Args[1]
	outputFile := "channels_followers.xlsx"

	// Initialize MongoDB client
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize components
	repo, err := database.NewMongoRepository(mongoClient, config)
	if err != nil {
		log.Fatalf("Error creating MongoDB repository: %v", err)
	}

	// Initialize components
	fm := filemanager.NewFileManager()
	telegramExtractor := telegram.NewTelegramExtractor()
	rutubeExtractor := rutube.NewRutubeExtractor()
	vkExtractor := vk.NewVKExtractor()
	instagramExtractor := instagram.NewInstagramExtractor()

	// Initialize and run app
	application := app.NewInfluencerApp(repo, fm, telegramExtractor, rutubeExtractor, vkExtractor, instagramExtractor)
	application.Run(inputFile, outputFile)

	log.Printf("Execution time: %v", time.Since(startAt))
}
