package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) EstimateTimeHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	inputFile := "uploaded_" + file.Filename
	if err := c.SaveFile(file, inputFile); err != nil {
		return err
	}

	// Estimate the processing time
	estimatedTime, err := h.InfluencerApp.EstimateProcessingTime(inputFile)
	if err != nil {
		return err
	}

	// Delete the temp file
	if err := os.Remove(inputFile); err != nil {
		log.Printf("Failed to delete temp file %s: %v", inputFile, err)
	}

	return c.JSON(fiber.Map{
		"estimatedTime": (estimatedTime) * 4,
	})
}
