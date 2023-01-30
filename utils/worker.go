// Worker - вспомогательный бекенд для различного дерьма, которое мне лень обрабатывать здесь
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"toncap-backend/types"
)

var worker_host = os.Getenv("WORKER_HOST")

func GetActual(contract string) (actual types.ActualResponse, err error) {
	response, err := http.Get(fmt.Sprintf("http://%s/markets/%s", worker_host, contract))
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

func GetCurrency() (types.Currency, error) {
	response, err := http.Get(fmt.Sprintf("http://%s/currency", worker_host))
	if err != nil {
		return types.Currency{}, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.Currency{}, err
	}
	defer response.Body.Close()

	var currency types.Currency
	err = json.Unmarshal(responseBodyRaw, &currency)
	if err != nil {
		return types.Currency{}, err
	}

	return currency, nil
}
