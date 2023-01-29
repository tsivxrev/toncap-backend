package router

import (
	"errors"
	"time"
	"toncap-backend/config"
	"toncap-backend/controllers"
	"toncap-backend/middlewares"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	cacheStore := persistence.NewInMemoryStore(time.Second)

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(middlewares.ErrorHandler))

	router.GET("/ads", controllers.GetAds)
	router.GET("/currency",
		cache.CachePage(
			cacheStore,
			config.CURRENCY_UPDATE_TIME,
			controllers.GetCurrency,
		),
	)

	router.Use(middlewares.Auth)

	router.GET("/jettons", controllers.GetJettons)
	router.GET("/jettons/:id", controllers.GetJetton)

	router.GET("/ads/_update", controllers.UpdateAds)
	router.GET("/jettons/_update", controllers.UpdateJettons)

	router.GET("/token/validate/:token", controllers.ValidateTokenController)
	router.GET("/token/generate", controllers.GenerateTokenController)

	router.POST("/prices", controllers.AddPrice)
	router.POST("/price", controllers.AddPrice)

	router.GET("/prices/graph/:contract", cache.CachePage(
		cacheStore,
		config.PRICE_UPDATE_TIME,
		controllers.GetGraph,
	))

	router.GET("/prices/min/:contract", cache.CachePage(
		cacheStore,
		config.PRICE_UPDATE_TIME,
		controllers.GetPriceMinimal,
	))

	router.GET("/prices/:contract",
		cache.CachePage(
			cacheStore,
			config.PRICE_UPDATE_TIME,
			controllers.GetPrice,
		),
	)

	router.NoMethod(func(c *gin.Context) {
		controllers.NewError(c, 405, errors.New("method is not allowed"))
	})

	router.NoRoute(func(c *gin.Context) {
		controllers.NewError(c, 404, errors.New("not found"))
	})

	return router
}
