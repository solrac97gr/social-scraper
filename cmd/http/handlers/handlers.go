package handlers

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/solrac97gr/telegram-followers-checker/app"
	"github.com/solrac97gr/telegram-followers-checker/database"
	"github.com/solrac97gr/telegram-followers-checker/extractors/extractor"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

type Handlers struct {
	Repository  database.InfluencerRepository
	FileManager filemanager.FileManager
	Extractors  []extractor.StatisticExtractor
}

func NewHandlers(repo database.InfluencerRepository, fm filemanager.FileManager, extractors ...extractor.StatisticExtractor) *Handlers {
	return &Handlers{
		Repository:  repo,
		FileManager: fm,
		Extractors:  extractors,
	}
}

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

	// Initialize and run app
	application := app.NewApp(h.Repository, h.FileManager, h.Extractors...)
	results := application.Run(inputFile, outputFile)

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

func (h *Handlers) DownloadHandler(c *fiber.Ctx) error {
	outputFile := c.Query("filename")
	return c.Download(outputFile)
}

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
	estimatedTime, err := h.FileManager.EstimateProcessingTime(inputFile)
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

func (h *Handlers) HealthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func (h *Handlers) AnalysesHandler(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}
	limit := c.Query("limit", "10")
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit number",
		})
	}
	// Ensure the limit does not exceed a reasonable maximum
	if limitNum > 100 {
		limitNum = 100
	}

	analyses, err := h.Repository.GetAllInfluencerAnalyses(pageNum, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve analyses",
		})
	}

	return c.JSON(analyses)
}
