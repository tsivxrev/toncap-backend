package controllers

import (
	"github.com/gin-gonic/gin"

	"toncap-backend/types"
)

func NewError(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, &types.HTTPError{
		Code:    status,
		Message: err.Error(),
	})
}
