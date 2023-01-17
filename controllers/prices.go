package controllers

import (
	"errors"
	"toncap-backend/database"
	"toncap-backend/logger"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func average(prices []types.Price) (price float64, volume float64) {
	for _, _price := range prices {
		price += _price.Price
		volume += _price.Volume
	}

	return price / float64(len(prices)), volume
}

func getPrice(contract string) gin.H {
	var prices []types.Price
	result := database.DB.Limit(14880).Find(&prices, "contract = ?", contract)
	logger.Log.Debugf("len: %v result: %v", len(prices), result)

	if len(prices) < 14880 {
		remaining := 14880 - len(prices)

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

	averages := map[string]map[string]float64{
		"actual": {
			"price":  0,
			"volume": 0,
		},
		"day": {
			"price":  0,
			"volume": 0,
		},
		"week": {
			"price":  0,
			"volume": 0,
		},
		"month": {
			"price":  0,
			"volume": 0,
		},
	}

	dayPrices := []types.Price{}   // len: 480
	weekPrices := []types.Price{}  // 3360
	monthPrices := []types.Price{} // 14880

	for idx, price := range prices {
		if idx == 0 {
			averages["actual"]["price"] += (price.Price / float64(len(prices)))
			averages["actual"]["volume"] += price.Volume
		}
		if idx+1 <= 480 {
			dayPrices = append(dayPrices, price)
		}
		if idx+1 <= 3360 {
			weekPrices = append(weekPrices, price)
		}
		if idx+1 <= 14880 {
			monthPrices = append(monthPrices, price)
		}
	}

	averages["day"]["price"], averages["day"]["volume"] = average(dayPrices)
	averages["week"]["price"], averages["week"]["volume"] = average(weekPrices)
	averages["month"]["price"], averages["month"]["volume"] = average(monthPrices)

	graph := make(map[int]map[string]float64)

	start := 0
	end := 479
	for i := 0; i <= 30; i++ {
		price, volume := average(monthPrices[start:end])
		graph[i] = map[string]float64{
			"price":  price,
			"volume": volume,
		}

		start += 480
		end += 480
	}

	return gin.H{
		"contract": contract,
		"averages": averages,
		"graph":    graph,
	}
}

func GetPrice(c *gin.Context) {
	contract := c.Param("contract")

	if contract == "" {
		NewError(c, 400, errors.New("contract not found"))
		return
	}

	c.JSON(200, getPrice(contract))
}

func AddPrice(c *gin.Context) {
	if c.GetString("auth_token_type") != "service" {
		NewError(c, 403, errors.New("access denied"))
		return
	}

	var price types.Price
	err := c.ShouldBindJSON(&price)
	if err != nil {
		NewError(c, 400, errors.New("invalid input"))
		return
	}

	result := database.DB.Create(&price)
	if result.Error != nil {
		logger.Log.Error(result.Error)
		NewError(c, 400, result.Error)
		return
	}

	logger.Log.Debug(price)
}
