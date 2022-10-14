package routes

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func PaymentRoute(version *gin.RouterGroup) {
	payment := version.Group("/payments")

	authPayment := payment.Use(api.AuthMiddleware())
	authPayment.GET("", api.RateLimit(10, time.Minute), api.GetPaymentHistory())
}
