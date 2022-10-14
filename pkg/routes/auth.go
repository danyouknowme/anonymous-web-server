package routes

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func AuthRoute(version *gin.RouterGroup) {
	auth := version.Group("/auth")
	{
		auth.POST("/register", api.RateLimit(10, time.Minute), api.CreateUser())
		auth.POST("/login", api.RateLimit(10, time.Minute), api.LoginUser())
	}
}
