package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{

		"status":  "ok",
		"message": "API TEST",
	})
}
