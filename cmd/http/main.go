package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/telegram-followers-checker/app"
	instagram "github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	vk "github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

func main() {
	fiberApp := fiber.New()

	fiberApp.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		inputFile := "uploaded_" + file.Filename
		if err := c.SaveFile(file, inputFile); err != nil {
			return err
		}

		outputFile := "channels_followers.xlsx"

		// Initialize components
		fm := filemanager.NewFileManager()
		telegramExtractor := telegram.NewTelegramExtractor()
		rutubeExtractor := rutube.NewRutubeExtractor()
		vkExtractor := vk.NewVKExtractor()
		instagramExtractor := instagram.NewInstagramExtractor()

		// Initialize and run app
		application := app.NewApp(fm, telegramExtractor, rutubeExtractor, vkExtractor, instagramExtractor)
		application.Run(inputFile, outputFile)

		// Download the output file
		if err := c.Download(outputFile); err != nil {
			return err
		}

		// Delete the temp files
		if err := os.Remove(inputFile); err != nil {
			log.Printf("Failed to delete temp file %s: %v", inputFile, err)
		}
		if err := os.Remove(outputFile); err != nil {
			log.Printf("Failed to delete temp file %s: %v", outputFile, err)
		}

		return nil
	})

	log.Fatal(fiberApp.Listen(":3000"))
}
