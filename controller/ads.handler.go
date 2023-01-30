package controller

import (
	"errors"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AdsGetAll(c *fiber.Ctx) error {
	var ads []types.Ad
	database.DB.Find(&ads)

	return c.Status(200).JSON(ads)
}

func AdsGetById(c *fiber.Ctx) error {
	id := c.Params("id")

	var ad types.Ad
	if database.DB.First(&ad, id).RowsAffected == 0 {
		return Error(c, fiber.StatusNotFound, errors.New("ad not found"))
	}

	return c.Status(200).JSON(ad)
}

func AdsUpdate(c *fiber.Ctx) error {
	// Все-таки стоит придумать более удобную проверку на тип токена
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	id := c.Params("id")

	var ad types.Ad
	if database.DB.First(&ad, id).RowsAffected == 0 {
		return Error(c, fiber.StatusNotFound, errors.New("ad not found"))
	}

	var updatedAd types.Ad
	err := c.BodyParser(&updatedAd)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	err = utils.Validate.Struct(updatedAd)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	updatedAd.Id = ad.Id //prevent id from changing

	database.DB.Save(&updatedAd)
	return c.Status(200).JSON(updatedAd)
}

func AdsCreate(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	var ad types.Ad
	err := c.BodyParser(&ad)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	err = utils.Validate.Struct(ad)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	ad.Id = 0
	database.DB.Create(&ad)

	return c.Status(fiber.StatusCreated).JSON(ad)
}

func AdsDelete(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	id := c.Params("id")

	var ad types.Ad
	if database.DB.First(&ad, id).RowsAffected == 0 {
		return Error(c, fiber.StatusNotFound, errors.New("ad not found"))
	}

	database.DB.Delete(ad)

	return c.Status(200).JSON(ad)
}
