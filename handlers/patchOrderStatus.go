package handlers

import (
	"AdminPanelAPI/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var OrderPatchBodyKey = "order_patch_body"

func PatchOrderStatus(c *gin.Context) {

	body, _ := c.Get(OrderPatchBodyKey)
	orderPayload := body.(models.OrderPatchPayload)

	c.JSON(http.StatusOK, gin.H{

		"valid":   true,
		"status":  "ok",
		"message": "response test",
		"paylaod": orderPayload,
	})
}
