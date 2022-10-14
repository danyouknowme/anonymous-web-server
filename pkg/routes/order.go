package routes

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func OrderRoute(version *gin.RouterGroup) {
	orders := version.Group("/orders")
	authOrders := orders.Use(api.AuthMiddleware())
	authOrders.POST("", api.RateLimit(10, time.Minute), api.AddOrder())

	authAndAdminOrders := orders.Use(api.AuthAndAdminMiddleWare())
	authAndAdminOrders.GET("", api.GetAllOrders())
	authAndAdminOrders.POST("/confirmation", api.ConfirmOrder())
}
