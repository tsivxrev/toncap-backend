package controller

import (
	"errors"
	"time"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func get_extended_graph(contract string) (graph map[string][]types.Graph) {
	var markets []string
	database.DB.Raw("SELECT DISTINCT market FROM prices WHERE contract = ?", contract).Scan(&markets)

	now := time.Now()

	graph = make(map[string][]types.Graph)
	for _, market := range markets {
		graph[market] = get_graph(contract, market, now)
	}

	graph["average"] = get_graph(contract, "", now)

	return graph
}

func get_graph(contract string, market string, date time.Time) (graph []types.Graph) {
	var prices []types.Price

	query := database.DB.Where("contract = ?", contract).Where("month = ?", int(date.Month())).Where("year = ?", date.Year())
	if market != "" {
		query = query.Where("market = ?", market)
	}

	query.Find(&prices)

	for _, price := range prices {
		price_date := time.Date(price.Year, time.Month(price.Month), price.Day, 23, 59, 0, 0, time.UTC)
		graph = append(graph, types.Graph{
			Date:   price_date.Unix(),
			Price:  price.Price,
			Volume: price.Volume,
		})
	}

	if market != "" {
		return graph
	}

	var markets []string
	database.DB.Raw("SELECT DISTINCT market FROM prices WHERE contract = ?", contract).Scan(&markets)

	var merged_graph []types.Graph
	temp_graph := make(map[int64]types.Graph)
	for _, graph_item := range graph {
		if val, ok := temp_graph[graph_item.Date]; ok {
			val.Price += graph_item.Price
			val.Volume += graph_item.Volume
			temp_graph[graph_item.Date] = val
		} else {
			temp_graph[graph_item.Date] = graph_item
		}
	}

	for k, v := range temp_graph {
		merged_graph = append(merged_graph, types.Graph{
			Date:   k,
			Price:  v.Price / float64(len(markets)),
			Volume: v.Volume / float64(len(markets)),
		})
	}

	return merged_graph
}

func GetPrice(c *fiber.Ctx) error {
	contract := c.Params("contract")

	graph := get_graph(contract, "", time.Now())
	actual, err := utils.GetActual(contract)
	if err != nil {
		return Error(c, 500, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"contract": contract,
		"actual":   actual.Actual,
		"markets":  actual.Markets,
		"graph":    graph,
	})
}

func GetGraph(c *fiber.Ctx) error {
	contract := c.Params("contract")
	graph := get_extended_graph(contract)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"contract": contract,
		"graph":    graph,
	})
}

func GetMinimalPrice(c *fiber.Ctx) error {
	contract := c.Params("contract")

	return c.Status(fiber.StatusOK).JSON(contract)
}

func AddPrice(c *fiber.Ctx) error {
	if c.Locals("auth.token_type") != "service" {
		return Error(c, fiber.StatusForbidden, errors.New("access denied"))
	}

	var received_price types.Price
	err := c.BodyParser(&received_price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, errors.New("invalid request body"))
	}

	err = utils.Validate.Struct(received_price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, err)
	}

	var price types.Price
	result := database.DB.Where("contract = ?", received_price.Contract).Where("market = ?", received_price.Market).Where("day = ?", received_price.Day).Where("month = ?", received_price.Month).Where("year = ?", received_price.Year).First(&price)
	if err != nil {
		return Error(c, fiber.StatusBadRequest, result.Error)
	}

	if result.RowsAffected == 0 {
		price = received_price
		database.DB.Create(&price)
	} else {
		price.Price = received_price.Price
		price.Volume = received_price.Volume
		database.DB.Save(&price)
	}

	return c.Status(fiber.StatusCreated).JSON(price)
}
