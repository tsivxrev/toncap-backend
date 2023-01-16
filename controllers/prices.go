package controllers

import (
	"errors"
	"toncap-backend/database"
	"toncap-backend/logger"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

/* func calcPriceByDuration(prices []types.Price, markets_count int, duration int) map[string]float64 {
	writes_count := duration * markets_count
	i := 0

	prices_resp := map[string]float64{
		"price":  0,
		"volume": 0,
	}
	for _, price := range prices {
		if i == writes_count {
			break
		}

		prices_resp["price"] += price.Price
		prices_resp["volume"] += price.Volume
		i++
	}

	return prices_resp
}

func getPrice_t(ticker string) (gin.H, error) {
	jetton, ok := config.JETTONS[ticker]
	if !ok {
		return gin.H{}, errors.New("ticker not found")
	}

	markets_count := len(jetton.Markets)

	var prices []types.Price
	result := database.DB.Raw("SELECT * FROM prices WHERE ticker = ?", ticker).Limit(14400 * markets_count).Scan(&prices)
	logger.Log.Debugf("result: %v", result)

	market_data := map[string]interface{}{
		"actual": calcPriceByDuration(prices, markets_count, 1),
		"day":    calcPriceByDuration(prices, markets_count, 480),
		"week":   calcPriceByDuration(prices, markets_count, 3360),
		"month":  calcPriceByDuration(prices, markets_count, 14400),
	}

	graph := make(map[int]map[string]float64)

	return gin.H{
		"market_data": market_data,
	}, nil
}

func getPrice(ticker string, markets_count int, writes int, offset_days int) map[string]float64 {
	offset := (markets_count * 480) * offset_days

	writes_get := writes * markets_count

	var prices []types.Price
	result := database.DB.Raw("SELECT * FROM prices WHERE ticker = ?", ticker).Limit(writes_get).Offset(offset).Scan(&prices)

	logger.Log.Debugf("result: %v", result)

	var price float64 = 0
	var volume float64 = 0

	for _, price_ := range prices {
		price += price_.Price / float64(len(prices))
		volume += price_.Volume
	}

	return map[string]float64{
		"price":  price,
		"volume": volume,
	}
}

func GetPrices(c *gin.Context) {
	tickerName := c.Query("ticker")

	jetton, ok := config.JETTONS[tickerName]
	if !ok {
		NewError(c, 404, errors.New("ticker not found"))
		return
	}

	markets := jetton.Markets
	markets_count := len(markets)

	market_data := map[string]interface{}{
		"actual": getPrice(tickerName, markets_count, 1, 0),
		"day":    getPrice(tickerName, markets_count, 480, 0),
		"week":   getPrice(tickerName, markets_count, 3360, 0),
		"month":  getPrice(tickerName, markets_count, 14400, 0),
	}

	market_prices := map[string]interface{}{"test": true}
	graph := make(map[int]map[string]float64)

	for i := 0; i <= 30; i++ {
		graph[i] = getPrice(tickerName, markets_count, 480, i)
	}

	c.JSON(200, gin.H{
		"market_data":   market_data,
		"market_prices": market_prices,
		"graph":         graph,
	})
} */

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
func EditPrice(c *gin.Context)   {}
func RemovePrice(c *gin.Context) {}
