package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func AuthRoute(version *gin.RouterGroup) {
	auth := version.Group("/auth")
	{
		auth.POST("/register", api.CreateUser())
		auth.POST("/login", api.LoginUser())
	}
}
