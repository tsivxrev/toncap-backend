package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"toncap-backend/database"
	"toncap-backend/logger"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func fetch_actual(contract string) (actual types.ActualResponse, err error) {
	response, err := http.Get("http://127.0.0.1:3001/markets/" + contract)
	if err != nil {
		return types.ActualResponse{}, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.ActualResponse{}, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(responseBodyRaw, &actual)
	if err != nil {
		return types.ActualResponse{}, err
	}

	return actual, nil
}

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

func getPrice(contract string) (priceResp map[string]interface{}, err error) {
	var prices []types.Price
	result := database.DB.Limit(14880).Find(&prices, "contract = ?", contract)
	logger.Log.Debugf("len: %v result: %v", len(prices), result)

	var markets []string
	result = database.DB.Raw("SELECT DISTINCT market FROM prices WHERE contract = ?", contract).Scan(&markets)
	logger.Log.Debugf("len: %v result: %v", len(markets), result)

	if len(prices) == 0 {
		return nil, errors.New("contract not found")
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

	dayPrices := []types.Price{}
	weekPrices := []types.Price{}
	monthPrices := []types.Price{}

	for idx, price := range prices {
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

	actual, err := fetch_actual(contract)
	if err != nil {
		return nil, err
	}

	averages["actual"] = map[string]float64{
		"price":  actual.Actual.Price,
		"volume": actual.Actual.Volume,
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

	return map[string]interface{}{
		"contract": contract,
		"averages": averages,
		"graph":    graph,
		"markets":  actual.Markets,
	}, nil
}

func GetPriceMinimal(c *gin.Context) {
	contract := c.Param("contract")

	if contract == "" {
		NewError(c, 400, errors.New("contract not found"))
		return
	}

	actual, err := fetch_actual(contract)
	if err != nil {
		NewError(c, 502, err)
		return
	}

	c.JSON(200, gin.H{
		"contract": contract,
		"actual": gin.H{
			"price":  actual.Actual.Price,
			"volume": actual.Actual.Volume,
		},
		"markets": actual.Markets,
	})
}

func GetPrice(c *gin.Context) {
	contract := c.Param("contract")

	if contract == "" {
		NewError(c, 400, errors.New("contract not found"))
		return
	}

	price, err := getPrice(contract)
	if err != nil {
		NewError(c, 400, err)
		return
	}

	c.JSON(200, price)
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
