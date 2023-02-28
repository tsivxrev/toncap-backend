package controller

import (
	"errors"
	"log"
	"time"
	"toncap-backend/database"
	"toncap-backend/types"
	"toncap-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func ListedContracts(c *fiber.Ctx) error {
	jettons, err := utils.GetJettons()
	if err != nil {
		return Error(c, 500, err)
	}

	log.Printf("jettons %v\n", jettons["data"])

	var parsed_contracts []fiber.Map
	contracts, ok := jettons["data"].([]interface{})
	if !ok {
		return Error(c, 500, errors.New("something went wrong"))
	}

	for _, contract := range contracts {
		contract := contract.(string)
		contract_data, err := get_contract(contract)
		if err != nil {
			continue
		}

		parsed_contracts = append(parsed_contracts, contract_data)
	}

	return c.Status(200).JSON(parsed_contracts)
}

func Contracts(c *fiber.Ctx) error {
	jettons, err := utils.GetJettons()
	if err != nil {
		return Error(c, 500, err)
	}

	return c.Status(200).JSON(jettons)
}

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

	actual, err := utils.GetActual(contract)
	if err != nil {
		return Error(c, 500, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"actual":  actual.Actual,
		"markets": actual.Markets,
	})
}

func GetGraph(c *fiber.Ctx) error {
	contract := c.Params("contract")
	graph := get_extended_graph(contract)

	return c.Status(fiber.StatusOK).JSON(graph)
}

func get_contract(contract string) (resp fiber.Map, err error) {
	if contract == "" {
		return nil, errors.New("contract not provided")
	}

	graph := get_extended_graph(contract)
	actual, err := utils.GetActual(contract)
	if err != nil {
		return nil, err
	}
	meta, err := utils.JettonMeta(contract)
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"contract": contract,
		"graph":    graph,
		"meta":     meta,
		"actual":   actual.Actual,
		"markets":  actual.Markets,
	}, nil
}

func GetContract(c *fiber.Ctx) error {
	contract := c.Params("contract")

	contract_data, err := get_contract(contract)
	if err != nil {
		return Error(c, 500, err)
	}

	return c.Status(200).JSON(contract_data)
}

func GetJettonMeta(c *fiber.Ctx) error {
	contract := c.Params("contract")

	if contract == "" {
		Error(c, 400, errors.New("contract not provided"))
	}

	meta, err := utils.JettonMeta(contract)
	if err != nil {
		Error(c, 500, err)
	}

	return c.Status(200).JSON(meta)
}
