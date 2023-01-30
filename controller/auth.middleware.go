package controller

import (
	"errors"
	"strconv"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("x-toncap-token")

	if token == "" {
		return Error(c, fiber.StatusUnauthorized, errors.New("token not provided"))
	}

	tokenData, valid := utils.ValidateToken(token)
	if !valid {
		return Error(c, fiber.StatusUnauthorized, errors.New("invalid token"))
	}

	c.Locals("auth.user_id", tokenData.UserId)
	c.Locals("auth.token_type", tokenData.Type)
	c.Locals("auth.token_id", tokenData.Id)

	c.Set("x-user-id", strconv.Itoa(tokenData.UserId))
	c.Set("x-token-id", tokenData.Id)

	return c.Next()
}
