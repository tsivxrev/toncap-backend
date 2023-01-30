package controller

import (
	"errors"
	"log"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func max(prices []types.Price) float64 {
	var m, t float64
	for _, e := range prices {
		if e.Price > t {
			t = e.Price
			m = t
		}
	}
	return m
}

func avg(prices []types.Price) (price float64, volume float64) {
	price = max(prices)
	volume = 0

	for _, price := range prices {
		volume += price.Volume / float64(len(prices))
	}

	return price, volume
}

func getPrice(contract string) (fiber.Map, error) {
	var prices []types.Price
	result := database.DB.Limit(14880).Find(&prices, "contract = ?", contract)
	if result.RowsAffected == 0 {
		return nil, errors.New("contract not found")
	}

	var markets []string
	result = database.DB.Raw("SELECT DISTINCT market FROM prices WHERE contract = ?", contract).Scan(&markets)
	if result.RowsAffected == 0 {
		return nil, errors.New("markets not found")
	}

	markets_count := len(markets)

	if len(prices) < markets_count*14880 {
		remaining := markets_count*14880 - len(prices)
		for i := 0; i < remaining; i++ {
			prices = append(prices, types.Price{
				Contract: prices[0].Contract,
				Ticker:   prices[0].Ticker,
				Market:   prices[0].Market,
				Volume:   0,
				Price:    0,
			})
		}
	}

	graph := make(map[int]map[string]float64)

	start := 0
	end := 479 * markets_count
	for i := 0; i <= 30; i++ {
		price, volume := avg(prices[start:end])
		graph[i] = map[string]float64{
			"price":  price,
			"volume": volume,
		}

		start += 480 * markets_count
		end += 480 * markets_count
	}

	actual, err := utils.GetActual(contract)
	if err != nil {
		return nil, nil
	}

	return fiber.Map{
		"averages": fiber.Map{
			"actual": fiber.Map{
				"price":  actual.Actual.Price,
				"volume": actual.Actual.Volume,
			},
		},
		"contract": contract,
		"graph":    graph,
		"markets":  markets,
	}, nil
}

func GetPrice(c *fiber.Ctx) error {
	contract := c.Params("contract")
	prices, err := getPrice(contract)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(prices)
}

func GetGraph(c *fiber.Ctx) error {
	contract := c.Params("contract")
	prices, err := getPrice(contract)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(prices["graph"])
}

func GetMinimalPrice(c *fiber.Ctx) error {
	contract := c.Params("contract")
	if contract == "" {
		return Error(c, fiber.StatusNotFound, errors.New("contract not found"))
	}

	actual, err := utils.GetActual(contract)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"contract": contract,
		"actual": fiber.Map{
			"price":  actual.Actual.Price,
			"volume": actual.Actual.Volume,
		},
		"markets": actual.Markets,
	})
}

func AddPrice(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	var price types.Price
	err := c.BodyParser(&price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	err = utils.Validate.Struct(price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	result := database.DB.Create(&price)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Printf("[AddPrice] %v result: %v\n", price, result)
		return Error(c, fiber.StatusInternalServerError, result.Error)
	}

	return c.Status(fiber.StatusCreated).JSON(price)
}
