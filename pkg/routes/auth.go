package routes

import (
	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
	router.POST("/register", api.CreateUser())
}
