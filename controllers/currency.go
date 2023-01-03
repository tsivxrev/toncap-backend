package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"toncap-backend/config"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

// TODO:
// https://stackoverflow.com/questions/18412126/golang-parse-a-json-with-dynamic-key
func GetCurrency(c *gin.Context) {
	response, err := http.Get("https://api.coingecko.com/api/v3/simple/token_price/ethereum?vs_currencies=usd&contract_addresses=" + config.CONTRACT_ADDRESS)
	if err != nil {
		NewError(c, 500, err)
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		NewError(c, 500, err)
	}
	defer response.Body.Close()

	tonPriceResponse := make(map[string]map[string]float64)
	err = json.Unmarshal(responseBodyRaw, &tonPriceResponse)
	if err != nil {
		NewError(c, 500, err)
	}

	response, err = http.Get("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/usd.json")
	if err != nil {
		NewError(c, 500, err)
	}

	responseBodyRaw, err = ioutil.ReadAll(response.Body)
	if err != nil {
		NewError(c, 500, err)
	}
	defer response.Body.Close()

	currenciesResponse := make(map[string]json.RawMessage)
	err = json.Unmarshal(responseBodyRaw, &currenciesResponse)
	if err != nil {
		NewError(c, 500, err)
	}

	tonPrice := tonPriceResponse[config.CONTRACT_ADDRESS]["usd"]
	currencies := make(map[string]float64)
	err = json.Unmarshal(currenciesResponse["usd"], &currencies)
	if err != nil {
		NewError(c, 500, err)
	}

	c.JSON(200, types.Currency{
		Byn: tonPrice * currencies["byn"],
		Cny: tonPrice * currencies["cny"],
		Eur: tonPrice * currencies["eur"],
		Gbp: tonPrice * currencies["gbp"],
		Kzt: tonPrice * currencies["kzt"],
		Rub: tonPrice * currencies["rub"],
		Uah: tonPrice * currencies["uah"],
		Usd: tonPrice * currencies["usd"],
	})
}
