package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/telegram-followers-checker/app"
)

type Handlers struct {
	InfluencerApp *app.InfluencerApp
	UsersApp      *app.UserApp
}

func NewHandlers(influencerApp *app.InfluencerApp, usersApp *app.UserApp) (*Handlers, error) {
	if influencerApp == nil || usersApp == nil {
		return nil, errors.New("app cannot be nil")
	}
	return &Handlers{
		InfluencerApp: influencerApp,
		UsersApp:      usersApp,
	}, nil
}

func (h *Handlers) HealthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}
