package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) LogoutUserHandler(c *fiber.Ctx) error {
	// Get user ID from context (set by JWT middleware)
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Invalidate user tokens
	err := h.UsersApp.LogoutUser(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
