package controller

import (
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Currency(c *fiber.Ctx) error {
	currency, err := utils.GetCurrency()
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(currency)
}
