package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func OrderRoute(version *gin.RouterGroup) {
	orders := version.Group("/orders")

	authAndAdminOrders := orders.Use(api.AuthAndAdminMiddleWare())
	authAndAdminOrders.GET("", api.GetAllOrders())
}
