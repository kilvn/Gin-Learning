package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(err error, c gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": err,
	})
}