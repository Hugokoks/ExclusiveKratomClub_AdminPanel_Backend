package handlers

import (
	"AdminPanelAPI/apperrors"
	"AdminPanelAPI/db"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrderByID(c *gin.Context){


	id := c.Params.ByName("id")

	ctx := c.Request.Context()

	order, err := db.SelectOrderDetail(ctx, id)
	if err != nil {
		log.Printf("!!! DATABASE ERROR: %v", err)

		if errors.Is(err,apperrors.ErrOrdersNotFound ){
			c.JSON(http.StatusNotFound, gin.H{ "valid":false,"message": "Order not found"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"valid":false,"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid":true,"message": "Order retrieved successfully","order": order})
}