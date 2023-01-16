package middlewares

import (
	"errors"
	"toncap-backend/controllers"
	"toncap-backend/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err any) {
	logger.Log.Error(err)
	controllers.NewError(c, 502, errors.New("oops! something went wrong"))
}
