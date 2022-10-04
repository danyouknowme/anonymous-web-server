package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func ResourceRoute(version *gin.RouterGroup) {
	resources := version.Group("/resources")
	{
		resources.GET("", api.GetAllResourcesInfo())
	}
}
