package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/solrac97gr/telegram-followers-checker/app"
	instagram "github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	vk "github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

func main() {
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	// Serve static files from the root directory
	fiberApp.Static("/", "./public")

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
		results := application.Run(inputFile, outputFile)

		// Delete the temp files
		if err := os.Remove(inputFile); err != nil {
			log.Printf("Failed to delete temp file %s: %v", inputFile, err)
		}

		return c.JSON(fiber.Map{
			"results": results,
		})
	})

	// Add a new endpoint to serve the results as a downloadable file
	fiberApp.Get("/download", func(c *fiber.Ctx) error {
		outputFile := "channels_followers.xlsx"
		return c.Download(outputFile)
	})

	log.Fatal(fiberApp.Listen(":3000"))
}
