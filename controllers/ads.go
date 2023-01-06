package controllers

import (
	"errors"
	"toncap-backend/config"
	"toncap-backend/types"

	"github.com/gin-gonic/gin"
)

func getAds() ([]types.Ad, error) {
	//	adsFile, err := os.ReadFile(config.ADS_FILE)
	//	if err != nil {
	//		return nil, nil
	//	}
	//
	//	var ads []types.Ad
	//	err = json.Unmarshal(adsFile, &ads)
	//	if err != nil {
	//		return nil, nil
	//	}
	//
	//	return ads, nil
	return config.ADS, nil
}

func UpdateAds(c *gin.Context) {
	if c.GetString("auth_token_type") == "service" {
		config.FetchAds()
		c.JSON(200, config.ADS)
		return
	}

	NewError(c, 403, errors.New("access denied"))
}

func GetAds(c *gin.Context) {
	ads, err := getAds()
	if err != nil {
		NewError(c, 500, err)
	}

	c.JSON(200, ads)
}

func GetAd(c *gin.Context) {
	id := c.Param("id")

	ads, err := getAds()
	if err != nil {
		NewError(c, 500, err)
	}

	for _, ad := range ads {
		if ad.Id == id {
			c.JSON(200, ad)
			return
		}
	}

	NewError(c, 404, errors.New("ad is not found"))
}

//func AddAd(c *gin.Context)    {}
//func EditAd(c *gin.Context)   {}
//func RemoveAd(c *gin.Context) {}
