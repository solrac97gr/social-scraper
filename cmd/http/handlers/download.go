package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func DownloadHandler(c *fiber.Ctx) error {
	outputFile := "channels_followers.xlsx"
	return c.Download(outputFile)
}
