package controller

import (
	"github.com/gofiber/fiber/v2"
)

func WebMiddleware(c *fiber.Ctx) error {
	//todo
	return c.Next()
}
