package controller

import (
	"errors"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GenerateToken(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	var tokenData types.TokenData
	err := c.BodyParser(&tokenData)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	err = utils.Validate.Struct(tokenData)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	tokenData.Id = uuid.NewString()

	token, err := utils.GenerateToken(tokenData)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(200).JSON(fiber.Map{
		"token": token,
		"data":  tokenData,
	})
}

func ValidateToken(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	token := c.Params("token")
	tokenData, valid := utils.ValidateToken(token)

	return c.Status(200).JSON(fiber.Map{
		"valid": valid,
		"data":  tokenData,
	})
}
