package router

import (
	"errors"
	"toncap-backend/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/currency", controller.Currency)
	app.Get("/records", controller.AuthMiddleware, controller.GetRecords)

	app.Get("/ads", controller.AdsGetAll)
	app.Get("/ads/:id", controller.AdsGetById)

	app.Post("/ads", controller.AuthMiddleware, controller.AdsCreate)
	app.Put("/ads/:id", controller.AuthMiddleware, controller.AdsUpdate)
	app.Delete("/ads/:id", controller.AuthMiddleware, controller.AdsDelete)

	app.Post("/prices", controller.AuthMiddleware, controller.AddPrice)

	app.Get("/contracts", controller.AuthMiddleware, controller.Contracts)
	app.Get("/contract/:contract", controller.AuthMiddleware, controller.GetContract)
	app.Get("/contract/:contract/meta", controller.AuthMiddleware, controller.GetJettonMeta)
	app.Get("/contract/:contract/graph", controller.AuthMiddleware, controller.GetGraph)
	app.Get("/contract/:contract/price", controller.AuthMiddleware, controller.GetPrice)

	app.Get("/token/generate", controller.AuthMiddleware, controller.GenerateToken)
	app.Get("/token/:token/validate", controller.AuthMiddleware, controller.ValidateToken)

	app.Use(func(c *fiber.Ctx) error {
		return controller.Error(c, fiber.StatusNotFound, errors.New("not found"))
	})
}
