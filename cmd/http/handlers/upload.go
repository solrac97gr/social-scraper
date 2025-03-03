package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/solrac97gr/telegram-followers-checker/app"
	"github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

func UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Ensure the directories exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll("results", os.ModePerm); err != nil {
		return err
	}

	uniqueID := uuid.New().String()
	inputFile := "uploads/" + uniqueID + "_uploaded_" + file.Filename
	if err := c.SaveFile(file, inputFile); err != nil {
		return err
	}

	outputFile := "results/" + uniqueID + "_channels_followers.xlsx"

	// Initialize components
	fm := filemanager.NewFileManager()
	telegramExtractor := telegram.NewTelegramExtractor()
	rutubeExtractor := rutube.NewRutubeExtractor()
	vkExtractor := vk.NewVKExtractor()
	instagramExtractor := instagram.NewInstagramExtractor()

	// Initialize and run app
	application := app.NewApp(fm, telegramExtractor, rutubeExtractor, vkExtractor, instagramExtractor)
	results := application.Run(inputFile, outputFile)

	// Delete the temp files
	if err := os.Remove(inputFile); err != nil {
		log.Printf("Failed to delete temp file %s: %v", inputFile, err)
	}

	return c.JSON(fiber.Map{
		"outputFile": outputFile,
		"results":    results,
	})
}
