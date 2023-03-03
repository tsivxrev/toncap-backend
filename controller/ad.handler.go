package controller //stable

import (
	"errors"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAds(c *fiber.Ctx) error {
	var ads []types.Ad
	query := database.DB.Find(&ads)
	if query.Error != nil {
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.Status(fiber.StatusOK).JSON(ads)
}

func GetAdById(c *fiber.Ctx) error {
	id := c.Params("id")

	var ad types.Ad
	query := database.DB.First(&ad, id)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return Error(c, fiber.StatusNotFound, errors.New("ad not found"))
		}
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.Status(fiber.StatusOK).JSON(ad)
}

func CreateAd(c *fiber.Ctx) error {
	var payload *types.AdCreateSchema
	err := c.BodyParser(&payload)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	validateErrors := utils.ValidateStruct(payload)
	if validateErrors != nil {
		return Error(c, fiber.StatusBadRequest, errors.New(utils.ValidateErrorString(validateErrors)))
	}

	newAd := &types.Ad{
		Type:     payload.Type,
		Text:     payload.Text,
		ImageURL: payload.ImageURL,
		Link:     payload.Link,
	}

	query := database.DB.Create(&newAd)
	if query.Error != nil {
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.Status(fiber.StatusCreated).JSON(newAd)
}

func UpdateAd(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload *types.AdCreateSchema
	err := c.BodyParser(&payload)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	validateErrors := utils.ValidateStruct(payload)
	if validateErrors != nil {
		return Error(c, fiber.StatusBadRequest, errors.New(utils.ValidateErrorString(validateErrors)))
	}

	var ad types.Ad
	query := database.DB.First(&ad, id)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return Error(c, fiber.StatusNotFound, errors.New("ad not found"))
		}
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	updates := make(map[string]any)
	if payload.Type != "" {
		updates["type"] = payload.Type
	}
	if payload.Text != "" {
		updates["text"] = payload.Text
	}
	if payload.ImageURL != "" {
		updates["image_url"] = payload.ImageURL
	}
	if payload.Link != "" {
		updates["link"] = payload.Link
	}

	query = database.DB.Model(&ad).Updates(updates)
	if query.Error != nil {
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.Status(fiber.StatusOK).JSON(ad)
}

func DeleteAd(c *fiber.Ctx) error {
	id := c.Params("id")

	query := database.DB.Delete(&types.Ad{}, "id = ?", id)

	if query.RowsAffected == 0 {
		return Error(c, fiber.StatusNotFound, errors.New("ad not found"))
	} else if query.Error != nil {
		return Error(c, fiber.StatusInternalServerError, query.Error)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
