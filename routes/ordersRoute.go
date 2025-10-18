package routes

import (
	"github.com/gin-gonic/gin"
	//ekc_mid "github.com/Hugokoks/kratomclub-go-common/middlewares"
	"AdminPanelAPI/handlers"
)

func OrdersRoute(g *gin.RouterGroup) {

	o := g.Group("/orders")

	o.GET("", handlers.GetOrders)

}
