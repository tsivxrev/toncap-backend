package controllers

import (
	"errors"
	"toncap-backend/database"
	"toncap-backend/logger"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
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

func getPrice(contract string) gin.H {
	var prices []types.Price
	result := database.DB.Limit(14880).Find(&prices, "contract = ?", contract)
	logger.Log.Debugf("len: %v result: %v", len(prices), result)

	var markets []string
	result = database.DB.Raw("SELECT DISTINCT market FROM prices WHERE contract = ?", contract).Scan(&markets)
	logger.Log.Debugf("len: %v result: %v", len(markets), result)

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
		if idx+1 <= 480*markets_count {
			dayPrices = append(dayPrices, price)
		}
		if idx+1 <= 3360*markets_count {
			weekPrices = append(weekPrices, price)
		}
		if idx+1 <= 14880*markets_count {
			monthPrices = append(monthPrices, price)
		}
	}

	averages["day"]["price"], averages["day"]["volume"] = avg(dayPrices)
	averages["week"]["price"], averages["week"]["volume"] = avg(weekPrices)
	averages["month"]["price"], averages["month"]["volume"] = avg(monthPrices)

	graph := make(map[int]map[string]float64)

	start := 0
	end := 479 * markets_count
	for i := 0; i <= 30; i++ {
		price, volume := avg(monthPrices[start:end])
		graph[i] = map[string]float64{
			"price":  price,
			"volume": volume,
		}

		start += 480 * markets_count
		end += 480 * markets_count
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
