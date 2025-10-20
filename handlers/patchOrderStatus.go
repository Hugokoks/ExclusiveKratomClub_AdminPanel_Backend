package handlers

import (
	"AdminPanelAPI/apperrors"
	"AdminPanelAPI/db"
	"AdminPanelAPI/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var OrderPatchBodyKey = "order_patch_body"

func PatchOrderStatus(c *gin.Context) {

	body, _ := c.Get(OrderPatchBodyKey)
	orderPayload := body.(models.OrderPatchPayload)
	if err := db.UpdateOrderStatus(c.Request.Context(), orderPayload.ID, orderPayload.Status); err != nil {

		if errors.Is(err, apperrors.ErrOrdersNotFound) {

			c.JSON(http.StatusNotFound, gin.H{"valid": false, "message": err.Error()})
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"valid": false, "message": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "Order status updated successfully.",
		"paylaod": orderPayload,
	})

}
