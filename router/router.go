package router

import (
	"errors"
	"toncap-backend/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/currency", controller.Currency)

	app.Get("/ads", controller.AdsGetAll)
	app.Get("/ads/:id", controller.AdsGetById)

	// Auth required routes
	app.Use(controller.AuthMiddleware)

	app.Post("/ads", controller.AdsCreate)
	app.Put("/ads/:id", controller.AdsUpdate)
	app.Delete("/ads/:id", controller.AdsDelete)

	app.Post("/prices", controller.AddPrice)
	app.Get("/prices/:contract", controller.GetPrice)
	app.Get("/prices/:contract/min", controller.GetMinimalPrice)
	app.Get("/prices/:contract/graph", controller.GetGraph)

	app.Get("/token/generate", controller.GenerateToken)
	app.Get("/token/:token/validate", controller.ValidateToken)

	app.Use(func(c *fiber.Ctx) error {
		return controller.Error(c, fiber.StatusNotFound, errors.New("not found"))
	})
}
