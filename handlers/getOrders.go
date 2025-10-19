package handlers

import (
	"AdminPanelAPI/apperrors"
	"AdminPanelAPI/db"
	"AdminPanelAPI/models"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {

	var filters models.OrderFilters

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid filter parameters"})
		return
	}

	orders, err := db.SelectOrders(c.Request.Context(), filters)
	log.Printf("!!! DATABASE ERROR: %v", err)

	if err != nil {
		if errors.Is(err, apperrors.ErrOrdersNotFound) {

			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No orders found."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch orders from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "orders retrive successfully",
		"orders":  orders,
	})
}
