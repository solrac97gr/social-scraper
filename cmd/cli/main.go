package main

import (
	"log"
	"os"
	"time"

	"github.com/solrac97gr/telegram-followers-checker/app"
	instagram "github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	vk "github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/config"
	"github.com/solrac97gr/telegram-followers-checker/integrations/tgstats/repository"
)

func main() {
	startAt := time.Now()

	// Check if input argument is provided
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the Excel file as an argument")
	}

	inputFile := os.Args[1]
	outputFile := "channels_followers.xlsx"

	// Initialize components
	fm := filemanager.NewFileManager()
	telegramExtractor := telegram.NewTelegramExtractor()
	rutubeExtractor := rutube.NewRutubeExtractor()
	vkExtractor := vk.NewVKExtractor()
	instagramExtractor := instagram.NewInstagramExtractor()

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

	// Initialize and run app
	application := app.NewApp(fm, repo, config, telegramExtractor, rutubeExtractor, vkExtractor, instagramExtractor)
	application.Run(inputFile, outputFile)

	log.Printf("Execution time: %v", time.Since(startAt))
}
