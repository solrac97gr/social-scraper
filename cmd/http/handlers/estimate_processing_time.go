package handlers

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) EstimateTimeHandler(c *fiber.Ctx) error {
	var inputFile string
	var needsCleanup bool

	// Check if it's text input or file upload
	textInput := c.FormValue("textInput")
	if textInput != "" {
		// Handle text input - create a temporary file
		uniqueID := "temp_" + strings.ReplaceAll(time.Now().Format("20060102150405"), " ", "_")
		inputFile = uniqueID + "_text_input.txt"
		if err := os.WriteFile(inputFile, []byte(textInput), 0644); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process text input",
			})
		}
		needsCleanup = true
	} else {
		// Handle file upload
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No file or text input provided",
			})
		}

		// Check file extension
		filename := strings.ToLower(file.Filename)
		if !strings.HasSuffix(filename, ".xlsx") &&
			!strings.HasSuffix(filename, ".xls") &&
			!strings.HasSuffix(filename, ".csv") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Unsupported file type. Please upload .xlsx, .xls, or .csv files",
			})
		}

		inputFile = "uploaded_" + file.Filename
		if err := c.SaveFile(file, inputFile); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save uploaded file",
			})
		}
		needsCleanup = true
	}

	// Use the existing EstimateProcessingTime method
	estimatedTime, err := h.InfluencerApp.EstimateProcessingTime(inputFile)
	if err != nil {
		// Clean up the temp file
		if needsCleanup {
			if err := os.Remove(inputFile); err != nil {
				log.Printf("Failed to delete temp file %s: %v", inputFile, err)
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to estimate processing time",
		})
	}

	// Delete the temp file
	if needsCleanup {
		if err := os.Remove(inputFile); err != nil {
			log.Printf("Failed to delete temp file %s: %v", inputFile, err)
		}
	}

	return c.JSON(fiber.Map{
		"estimatedTime": estimatedTime * 4,
	})
}
