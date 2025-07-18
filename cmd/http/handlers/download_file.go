package handlers

import "github.com/gofiber/fiber/v2"

func (h *Handlers) DownloadHandler(c *fiber.Ctx) error {
	outputFile := c.Query("filename")
	return c.Download(outputFile)
}
