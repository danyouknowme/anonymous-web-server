package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func ResourceRoute(router *gin.Engine) {
	router.GET("/resources", api.GetAllResourcesInfo())
}
