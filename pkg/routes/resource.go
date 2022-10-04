package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func ResourceRoute(version *gin.RouterGroup) {
	resources := version.Group("/resources")
	resources.GET("", api.GetAllResourcesInfo())
	resources.GET("/:resourceName", api.GetResourceByName())

	authResources := resources.Use(api.AuthMiddleware())
	authResources.GET("/download/:resourceName", api.GetDownloadResource())

	authAndAdminResources := resources.Use(api.AuthAndAdminMiddleWare())
	authAndAdminResources.POST("", api.CreateNewResource())
	authAndAdminResources.PATCH("", api.UpdateResource())
	authAndAdminResources.GET("/admin/:resourceName", api.GetResourceByNameAndRequiredAdmin())
}
