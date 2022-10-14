package routes

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/gin-gonic/gin"
)

func DownloadRoute(version *gin.RouterGroup) {
	download := version.Group("/downloads")
	authDownload := download.Use(api.AuthMiddleware())
	authDownload.GET("/:resourceName", api.RateLimit(5, time.Minute), api.GetDownloadResource())
	authDownload.POST("", api.DownloadResource())
}
