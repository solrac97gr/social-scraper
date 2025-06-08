package handlers

import "github.com/gofiber/fiber/v2"

func (h *Handlers) RegisterUserHandler(c *fiber.Ctx) error {
	// Parse the request body to get user details
	var userDetails struct {
		Username             string `json:"username"`
		Email                string `json:"email"`
		Password             string `json:"password"`
		ConfirmationPassword string `json:"confirmation_password"`
	}

	if err := c.BodyParser(&userDetails); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Call the UsersApp to register the user
	if err := h.UsersApp.SaveUser(userDetails.Username, userDetails.Email, userDetails.Password, userDetails.ConfirmationPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
}
