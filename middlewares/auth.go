package middlewares

import (
	"errors"
	"fmt"
	"toncap-backend/controllers"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	tokenString := c.GetHeader("x-toncap-token")

	if len(tokenString) < 1 {
		controllers.NewError(c, 401, errors.New("token not provided"))
		return
	}

	tokenData, valid := controllers.ValidateToken(tokenString)
	if !valid {
		controllers.NewError(c, 401, errors.New("invalid token"))
		return
	}

	c.Set("auth_user_id", tokenData.UserId)
	c.Set("auth_token_type", tokenData.Type)
	c.Header("x-user-id", fmt.Sprintf("%v", tokenData.UserId))

	c.Next()
}
