package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func OrderRoute(version *gin.RouterGroup) {
	orders := version.Group("/orders")
	authOrders := orders.Use(api.AuthMiddleware())
	authOrders.POST("", api.AddOrder())

	authAndAdminOrders := orders.Use(api.AuthAndAdminMiddleWare())
	authAndAdminOrders.GET("", api.GetAllOrders())
}
