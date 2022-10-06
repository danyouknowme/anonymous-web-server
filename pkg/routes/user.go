package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func UserRoute(version *gin.RouterGroup) {
	user := version.Group("/user")

	user.POST("/license", api.CheckLicense())
}
