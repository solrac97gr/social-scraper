package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

func EstimateTimeHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	inputFile := "uploaded_" + file.Filename
	if err := c.SaveFile(file, inputFile); err != nil {
		return err
	}

	// Initialize file manager
	fm := filemanager.NewFileManager()

	// Estimate the processing time
	estimatedTime, err := fm.EstimateProcessingTime(inputFile)
	if err != nil {
		return err
	}

	// Delete the temp file
	if err := os.Remove(inputFile); err != nil {
		log.Printf("Failed to delete temp file %s: %v", inputFile, err)
	}

	return c.JSON(fiber.Map{
		"estimatedTime": estimatedTime,
	})
}
