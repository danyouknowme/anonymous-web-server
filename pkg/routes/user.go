package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func UserRoute(version *gin.RouterGroup) {
	user := version.Group("/users")
	user.POST("/license", api.CheckLicense())

	authUsers := user.Use(api.AuthMiddleware())
	authUsers.POST("/ip/reset", api.ResetIP())

	authAndAdminUsers := user.Use(api.AuthAndAdminMiddleWare())
	authAndAdminUsers.GET("/:username", api.GetUserData())
	authAndAdminUsers.PATCH("", api.UpdateUserData())
}
