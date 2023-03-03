package router

import (
	"errors"
	"toncap-backend/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/currency", controller.GetCurrency)

	app.Get("/records", controller.AuthMiddleware, controller.GetRecords)
	app.Post("/records", controller.AuthMiddleware, controller.AuthServiceMethod, controller.CreateRecord)

	app.Get("/ads", controller.GetAds)
	app.Get("/ads/:id", controller.GetAdById)
	app.Post("/ads", controller.AuthMiddleware, controller.AuthServiceMethod, controller.CreateAd)
	app.Put("/ads/:id", controller.AuthMiddleware, controller.AuthServiceMethod, controller.UpdateAd)
	app.Delete("/ads/:id", controller.AuthMiddleware, controller.AuthServiceMethod, controller.DeleteAd)

	app.Get("/contracts", controller.GetContracts)
	app.Get("/contracts/listed", controller.GetListedContracts)
	app.Get("/contract/:contract", controller.GetContract)
	app.Get("/contract/:contract/meta", controller.GetContractMeta)
	app.Get("/contract/:contract/graph", controller.GetContractGraph)
	app.Get("/contract/:contract/price", controller.GetContractPrice)

	app.Get("/token/generate", controller.AuthMiddleware, controller.AuthServiceMethod, controller.GenerateToken)
	app.Get("/token/:token/validate", controller.AuthMiddleware, controller.AuthServiceMethod, controller.ValidateToken)

	app.Use(func(c *fiber.Ctx) error {
		return controller.Error(c, fiber.StatusNotFound, errors.New("not found"))
	})
}
