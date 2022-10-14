package routes

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func BenefitRoute(version *gin.RouterGroup) {
	benefit := version.Group("/benefits")
	authAndAdminBenefit := benefit.Use(api.AuthAndAdminMiddleWare())
	authAndAdminBenefit.GET("/", api.GetPartnerBenefits())
	authAndAdminBenefit.POST("/clear", api.RateLimit(10, 24*time.Hour), api.ClearPartnerBenefit())
}
