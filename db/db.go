package db

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"toncap-backend/storage"
	"toncap-backend/types"
)

func GetAds() ([]types.Ad, error) {
	adsFile, err := os.ReadFile(storage.App.AdsFilePath)
	if err != nil {
		return nil, err
	}

	var ads []types.Ad
	err = json.Unmarshal(adsFile, &ads)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func GetExchangeRate() (types.ExchangeRate, error) {
	response, err := http.Get("https://api.coingecko.com/api/v3/simple/token_price/ethereum?vs_currencies=usd&contract_addresses=" + storage.App.ContractAddress)
	if err != nil {
		return types.ExchangeRate{}, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.ExchangeRate{}, err
	}
	defer response.Body.Close()

	tonPriceResponse := make(map[string]map[string]float64)
	err = json.Unmarshal(responseBodyRaw, &tonPriceResponse)
	if err != nil {
		return types.ExchangeRate{}, err
	}

	response, err = http.Get("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/usd.json")
	if err != nil {
		return types.ExchangeRate{}, err
	}

	responseBodyRaw, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return types.ExchangeRate{}, err
	}
	defer response.Body.Close()

	currenciesResponse := make(map[string]json.RawMessage)
	err = json.Unmarshal(responseBodyRaw, &currenciesResponse)
	if err != nil {
		return types.ExchangeRate{}, err
	}

	tonPrice := tonPriceResponse[storage.App.ContractAddress]["usd"]
	currencies := make(map[string]float64)
	err = json.Unmarshal(currenciesResponse["usd"], &currencies)
	if err != nil {
		return types.ExchangeRate{}, err
	}

	exchangeRate := types.ExchangeRate{
		Byn: tonPrice * currencies["byn"],
		Cny: tonPrice * currencies["cny"],
		Eur: tonPrice * currencies["eur"],
		Gbp: tonPrice * currencies["gbp"],
		Kzt: tonPrice * currencies["kzt"],
		Rub: tonPrice * currencies["rub"],
		Uah: tonPrice * currencies["uah"],
		Usd: tonPrice * currencies["usd"],
	}

	return exchangeRate, nil
}
