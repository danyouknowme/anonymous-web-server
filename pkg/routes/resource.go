package routes

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func ResourceRoute(version *gin.RouterGroup) {
	resources := version.Group("/resources")
	resources.GET("", api.GetAllResourcesInfo())
	resources.GET("/:resourceName", api.GetResourceByName())

	authAndAdminResources := resources.Use(api.AuthAndAdminMiddleWare())
	authAndAdminResources.POST("", api.RateLimit(10, time.Minute), api.CreateNewResource())
	authAndAdminResources.PATCH("", api.UpdateResource())
	authAndAdminResources.GET("/admin/:resourceName", api.GetResourceByNameAndRequiredAdmin())
	authAndAdminResources.POST("/user", api.AddResourceToUser())
}
