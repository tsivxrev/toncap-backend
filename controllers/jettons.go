package controllers

import (
	"errors"
	"toncap-backend/config"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func getJettons() (map[string]types.Jetton, error) {
	//	jettonsFile, err := os.ReadFile(config.JETTONS_FILE)
	//	if err != nil {
	//		return nil, nil
	//	}
	//
	//	var jettons map[string]types.Jetton
	//	err = json.Unmarshal(jettonsFile, &jettons)
	//	if err != nil {
	//		return nil, nil
	//	}
	//
	//	return jettons, nil
	return config.JETTONS, nil
}

func UpdateJettons(c *gin.Context) {
	if c.GetString("auth_token_type") == "service" {
		config.FetchJettons()
		c.JSON(200, config.JETTONS)
		return
	}

	NewError(c, 403, errors.New("access denied"))
}

func GetJettons(c *gin.Context) {
	jettons, err := getJettons()
	if err != nil {
		NewError(c, 500, err)
	}

	keys := make([]string, 0, len(jettons))
	for k := range jettons {
		keys = append(keys, k)
	}

	c.JSON(200, gin.H{
		"all": keys,
		"tops": gin.H{
			"new":     []string{"fnz", "scale", "ambr"},
			"popular": []string{"fnz", "grbs", "hedge", "ambr", "take", "bolt"},
		},
	})
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
