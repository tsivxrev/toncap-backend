package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func GetCurrency(c *gin.Context) {
	response, err := http.Get("http://127.0.0.1:3001/currency")
	if err != nil {
		NewError(c, 500, err)
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		NewError(c, 500, err)
	}
	defer response.Body.Close()

	var currency types.Currency
	err = json.Unmarshal(responseBodyRaw, &currency)
	if err != nil {
		NewError(c, 500, err)
	}

	c.JSON(200, currency)
}
