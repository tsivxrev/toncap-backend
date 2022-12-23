package router

import (
	"toncap-backend/controllers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.New()

	router.GET("/ads", controllers.Ads)
	router.GET("/ton", controllers.ExchangeRate)

	return router
}
