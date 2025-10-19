package routes

import (
	"AdminPanelAPI/handlers"
	"AdminPanelAPI/limiter"

	ekc_mid "github.com/Hugokoks/kratomclub-go-common/middlewares"
	"github.com/gin-gonic/gin"
)

func OrdersRoute(g *gin.RouterGroup) {

	o := g.Group("/orders")

	o.GET("", ekc_mid.ValidateParamLength(limiter.BlogParamLimit, "query"), handlers.GetOrders)
	o.PATCH("/patch", handlers.PatchOrderStatus)

}
