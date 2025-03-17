package main

import (
	"fmt"

	"math/rand"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/stat", func(c *fiber.Ctx) error {
		token := c.Query("token")
		chennelId := c.Query("channelId")
		fmt.Println("Mock API recived:", token, chennelId)

		return c.JSON(fiber.Map{
			"status": "ok",
			"response": fiber.Map{
				"avg_post_reach": rand.Intn(100),
				"er_percent":     rand.Intn(100),
			},
		})
	})

	app.Listen(":3500")
}
