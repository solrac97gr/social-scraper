package main

import (
	"log"
	"os"

	"github.com/solrac97gr/telegram-followers-checker/app"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	vk "github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

func main() {
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

	// Initialize and run app
	application := app.NewApp(fm, telegramExtractor, rutubeExtractor, vkExtractor)
	application.Run(inputFile, outputFile)
}
