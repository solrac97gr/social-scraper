package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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

	analyses, err := h.InfluencerApp.GetAllInfluencerAnalysis(pageNum, limitNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve analyses",
		})
	}

	return c.JSON(analyses)
}
