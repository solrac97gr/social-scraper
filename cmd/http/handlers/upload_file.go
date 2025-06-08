package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handlers) UploadHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	// Ensure the directories exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create uploads directory",
		})
	}
	if err := os.MkdirAll("results", os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create results directory",
		})
	}

	uniqueID := uuid.New().String()
	outputFile := "results/" + uniqueID + "_channels_followers.xlsx"

	var inputFile string
	var needsCleanup bool

	// Check if this is text input or file upload
	textInput := c.FormValue("textInput")
	if textInput != "" {
		// For text input, create a temporary file to maintain compatibility
		inputFile = "uploads/" + uniqueID + "_text_input.txt"
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

		// Validate file type
		filename := strings.ToLower(file.Filename)
		if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") && !strings.HasSuffix(filename, ".csv") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Unsupported file format. Please upload .xlsx, .xls, or .csv files",
			})
		}

		inputFile = "uploads/" + uniqueID + "_uploaded_" + file.Filename
		if err := c.SaveFile(file, inputFile); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save uploaded file",
			})
		}
		needsCleanup = true
	}

	// Use the original Run method - FileManager will handle file type detection
	results := h.InfluencerApp.Run(userID, inputFile, outputFile)

	// Clean up temp file if needed
	if needsCleanup {
		if err := os.Remove(inputFile); err != nil {
			log.Printf("Failed to delete temp file %s: %v", inputFile, err)
		}
	}

	log.Printf("Processing completed for input: %s", inputFile)

	return c.JSON(fiber.Map{
		"outputFile": outputFile,
		"results":    results,
	})
}
