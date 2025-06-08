package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) LoginUserHandler(c *fiber.Ctx) error {
	// Extract email and password from the request body
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if credentials.Email == "" || credentials.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Authenticate the user and generate token
	token, user, err := h.UsersApp.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return token and user information (excluding sensitive data)
	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":               user.ID,
			"username":         user.Username,
			"email":            user.Email,
			"role":             user.Role,
			"subscription":     user.Subscription,
			"profile_complete": user.ProfileComplete,
			"created_at":       user.CreatedAt,
			"updated_at":       user.UpdatedAt,
		},
	})
}
