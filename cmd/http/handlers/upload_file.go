package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handlers) UploadHandler(c *fiber.Ctx) error {
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
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	// Run the application logic
	results := h.InfluencerApp.Run(userID, inputFile, outputFile)

	// Delete the temp files
	if err := os.Remove(inputFile); err != nil {
		log.Printf("Failed to delete temp file %s: %v", inputFile, err)
	}

	log.Printf("Processing completed for file: %s", inputFile)

	return c.JSON(fiber.Map{
		"outputFile": outputFile,
		"results":    results,
	})
}
