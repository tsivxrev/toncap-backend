package controllers

import (
	"io/ioutil"
	"net/http"

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
		return
	}
	defer response.Body.Close()

	c.Data(200, "application/json", responseBodyRaw)
}
