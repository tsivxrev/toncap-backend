package controller

import (
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(types.HTTPError{
		Code:    statusCode,
		Message: utils.Capitalize(err.Error()),
	})
}
