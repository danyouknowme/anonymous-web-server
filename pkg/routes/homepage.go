package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func HomepageRoute(router *gin.Engine) {
	router.GET("/homepage", api.GetHomepageInformation())
}
