package controller //stable

import (
	"errors"
	"strconv"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func GetRecords(c *fiber.Ctx) error {
	contract := c.Query("contract")
	limit_str := c.Query("limit", "10")

	limit, _ := strconv.Atoi(limit_str)

	query := database.DB
	if contract != "" {
		query = query.Where("contract = ?", contract)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	var prices []types.Price
	query = query.Find(&prices)
	if query.Error != nil {
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.Status(200).JSON(prices)
}

func CreateRecord(c *fiber.Ctx) error {
	var payload *types.CreatePriceSchema
	err := c.BodyParser(&payload)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	validateErrors := utils.ValidateStruct(payload)
	if validateErrors != nil {
		return Error(c, fiber.StatusBadRequest, errors.New(utils.ValidateErrorString(validateErrors)))
	}

	newPrice := &types.Price{
		Contract: payload.Contract,
		Ticker:   payload.Ticker,
		Market:   payload.Market,
		Price:    payload.Price,
		Volume:   payload.Volume,
		Date:     payload.Date,
	}
	query := database.DB.Create(&newPrice)
	if query.Error != nil {
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.Status(fiber.StatusCreated).JSON(newPrice)
}
