package controller //stable

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

func AuthUserMethod(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "user" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	return c.Next()
}

func AuthServiceMethod(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	return c.Next()
}
