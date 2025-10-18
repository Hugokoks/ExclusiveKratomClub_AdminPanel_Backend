package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderFilters struct {
	ID             string `form:"id"`
	FirstName      string `form:"firstName"`
	LastName       string `form:"lastName"`
	Email          string `form:"email"`
	Address        string `form:"address"`
	PaymentMethod  string `form:"paymentMethod"`
	DeliveryMethod string `form:"deliveryMethod"`
	Status         string `form:"status"`
	SortBy         string `form:"sortBy"`
	SortOrder      string `form:"sortOrder"`
}

func GetOrders(c *gin.Context) {

	var filters OrderFilters

	if err := c.ShouldBindQuery(&filters); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"status":  "ok",
		"message": "API TEST",
		"filters": filters,
	})
}
