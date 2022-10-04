package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func HomepageRoute(version *gin.RouterGroup) {
	homepage := version.Group("/homepage")
	homepage.GET("", api.GetHomepageInformation())

	authAndAdminHomepage := homepage.Use(api.AuthAndAdminMiddleWare())
	authAndAdminHomepage.PATCH("", api.UpdateHomepageInformation())
}
