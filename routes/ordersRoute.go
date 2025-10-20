package routes

import (
	"AdminPanelAPI/handlers"
	"AdminPanelAPI/limiter"
	"AdminPanelAPI/models"

	ekc_mid "github.com/Hugokoks/kratomclub-go-common/middlewares"
	"github.com/gin-gonic/gin"
)

func OrdersRoute(g *gin.RouterGroup) {

	o := g.Group("/orders")

	o.GET("", ekc_mid.ValidateParamLength(limiter.OrderFilterLimit, "query"), handlers.GetOrders)
	o.PATCH("/patch", ekc_mid.ValidateJSONBody(handlers.OrderPatchBodyKey, func(payload *models.OrderPatchPayload) error { return nil }),
		handlers.PatchOrderStatus)

}
