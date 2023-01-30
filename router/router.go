package router

import (
	"errors"
	"toncap-backend/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return controller.Error(c, fiber.StatusNotFound, errors.New("not found"))
	})
}
