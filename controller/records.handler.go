package controller

import (
	"strconv"
	"toncap-backend/database"
	"toncap-backend/types"

	"github.com/gofiber/fiber/v2"
)

func GetRecords(c *fiber.Ctx) error {
	contract := c.Query("contract")
	limit, err := strconv.Atoi(c.Query("limit"))

	var prices []types.Price
	query := database.DB.Where("contract = ?", contract)

	if limit > 0 || err == nil {
		query = query.Limit(limit)
	}

	result := query.Find(&prices)
	if result.Error != nil {
		Error(c, 500, result.Error)
	}

	return c.Status(200).JSON(prices)
}
