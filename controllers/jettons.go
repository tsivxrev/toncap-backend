package controllers

import (
	"encoding/json"
	"errors"
	"os"
	"toncap-backend/config"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func getJettons() (map[string]types.Jetton, error) {
	jettonsFile, err := os.ReadFile(config.JETTONS_FILE)
	if err != nil {
		return nil, nil
	}

	var jettons map[string]types.Jetton
	err = json.Unmarshal(jettonsFile, &jettons)
	if err != nil {
		return nil, nil
	}

	return jettons, nil
}

func GetJettons(c *gin.Context) {
	jettons, err := getJettons()
	if err != nil {
		NewError(c, 500, err)
	}

	c.JSON(200, jettons)
}

func GetJetton(c *gin.Context) {
	id := c.Param("id")

	jettons, err := getJettons()
	if err != nil {
		NewError(c, 500, err)
	}

	jetton, ok := jettons[id]
	if ok {
		c.JSON(200, jetton)
		return
	}

	NewError(c, 404, errors.New("jetton is not found"))
}

//func AddJetton(c *gin.Context)    {}
//func EditJetton(c *gin.Context)   {}
//func RemoveJetton(c *gin.Context) {}
