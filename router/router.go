package router

import (
	"errors"
	"time"
	"toncap-backend/config"
	"toncap-backend/controllers"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true

	cacheStore := persistence.NewInMemoryStore(time.Second)

	api := router.Group("/api/")

	api.GET("/ads", cache.CachePage(cacheStore, config.ADS_UPDATE_TIME, controllers.GetAds))
	//api.GET("/ads/:id", cache.CachePage(cacheStore, config.ADS_UPDATE_TIME, controllers.GetAd))
	//api.POST("/ads", controllers.AddAd)
	//api.PATCH("/ads/:id", controllers.EditAd)
	//api.DELETE("/ads/:id", controllers.RemoveAd)

	api.GET("/prices", controllers.GetPrices)
	api.GET("/prices/:id", controllers.GetPrice)
	api.POST("/prices", controllers.AddPrice)
	api.PATCH("/prices", controllers.EditPrice)
	api.DELETE("/prices", controllers.RemovePrice)

	api.GET("/currency", cache.CachePage(cacheStore, config.CURRENCY_UPDATE_TIME, controllers.GetCurrency))

	api.GET("/jettons", cache.CachePage(cacheStore, config.JETTONS_UPDATE_TIME, controllers.GetJettons))
	api.GET("/jettons/:id", cache.CachePage(cacheStore, config.JETTONS_UPDATE_TIME, controllers.GetJetton))
	//api.POST("/jettons/:id", controllers.AddJetton)
	//api.PATCH("/jettons/:id", controllers.EditJetton)
	//api.DELETE("/jettons/:id", controllers.RemoveJetton)

	router.NoMethod(func(c *gin.Context) {
		controllers.NewError(c, 400, errors.New("method is not allowed"))
	})

	router.NoRoute(func(c *gin.Context) {
		controllers.NewError(c, 404, errors.New("not found"))
	})

	return router
}
