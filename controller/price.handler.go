package controller

import (
	"errors"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AddPrice(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	var received_price types.Price
	err := c.BodyParser(&received_price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, errors.New("invalid request body"))
	}

	err = utils.Validate.Struct(received_price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	var price types.Price
	result := database.DB.Where("contract = ?", received_price.Contract).Where("market = ?", received_price.Market).Where("day = ?", received_price.Day).Where("month = ?", received_price.Month).Where("year = ?", received_price.Year).First(&price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, result.Error)
	}

	if result.RowsAffected == 0 {
		price = received_price
		database.DB.Create(&price)
	} else {
		price.Price = received_price.Price
		price.Volume = received_price.Volume
		database.DB.Save(&price)
	}

	return c.Status(fiber.StatusCreated).JSON(price)
}
