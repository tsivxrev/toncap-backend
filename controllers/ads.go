package controllers

import (
	"net/http"
	"toncap-backend/storage"

	"github.com/gin-gonic/gin"
)

func Ads(c *gin.Context) {
	c.JSON(http.StatusOK, storage.App.Ads)
}
