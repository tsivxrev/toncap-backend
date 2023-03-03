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

func getContract(contract string) (contractInfo types.ContractResponse, err error) {
	now := time.Now()
	from := now.AddDate(0, -1, 0).Unix()
	to := now.Unix()

	graph, err := calculateGraph(contract, "", from, to)
	if err != nil {
		log.Println(err)
	}

	actual, err := utils.GetActual(contract)
	if err != nil {
		return types.ContractResponse{}, err
	}

	meta, err := utils.GetContractMeta(contract)
	if err != nil {
		return types.ContractResponse{}, err
	}

	return types.ContractResponse{
		Contract: contract,
		Actual:   actual.Price,
		Markets:  actual.Markets,
		Meta:     meta,
		Graph:    graph,
	}, nil
}

func calculateGraph(contract string, market string, from int64, to int64) (graph []types.ContractGraph, err error) {
	var graphResult []types.ContractGraph

	query := database.DB.Model(&types.Price{})
	query = query.Select("date, AVG(price) AS price, SUM(volume) AS volume")
	query = query.Where("contract = ?", contract)
	if market != "" {
		query = query.Where("market = ?", market)
	}
	if from != 0 && to != 0 {
		query = query.Where("date BETWEEN ? AND ?", from, to)
	}
	query = query.Group("DATE(date, 'unixepoch')")
	query = query.Find(&graphResult)
	if query.Error != nil {
		return []types.ContractGraph{}, query.Error
	}

	return graphResult, nil
}

func GetListedContracts(c *fiber.Ctx) error {
	listedContracts, err := utils.GetListedContracts()
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	var contractsInfo []types.ContractResponse
	for _, contract := range listedContracts {
		contractInfo, err := getContract(contract)
		if err != nil {
			log.Println(err)
			continue
		}

		contractsInfo = append(contractsInfo, contractInfo)
	}

	return c.Status(fiber.StatusOK).JSON(contractsInfo)
}

func GetContractPrice(c *fiber.Ctx) error {
	contract := c.Params("contract")

	if contract == "" {
		return Error(c, fiber.StatusBadRequest, errors.New("contract not provided"))
	}

	actual, err := utils.GetActual(contract)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(actual)
}

func GetContractMeta(c *fiber.Ctx) error {
	contract := c.Params("contract")

	if contract == "" {
		return Error(c, fiber.StatusBadRequest, errors.New("contract not provided"))
	}

	meta, err := utils.GetContractMeta(contract)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(meta)
}

func GetContractGraph(c *fiber.Ctx) error {
	contract := c.Params("contract")

	if contract == "" {
		return Error(c, fiber.StatusBadRequest, errors.New("contract not provided"))
	}

	now := time.Now()
	from := now.AddDate(0, -1, 0).Unix()
	to := now.Unix()

	meta, err := calculateGraph(contract, "", from, to)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(meta)
}

func GetContracts(c *fiber.Ctx) error {
	listedContracts, err := utils.GetListedContracts()
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(listedContracts)
}

func GetContract(c *fiber.Ctx) error {
	contract := c.Params("contract")

	if contract == "" {
		return Error(c, fiber.StatusBadRequest, errors.New("contract not provided"))
	}

	contractInfo, err := getContract(contract)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(contractInfo)
}
