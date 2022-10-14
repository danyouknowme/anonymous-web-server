package main

import (
	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/routes"
	"github.com/danyouknowme/awayfromus/pkg/time"
	"github.com/danyouknowme/awayfromus/pkg/util"

	_ "github.com/danyouknowme/awayfromus/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Awayfromus API
// @version 1.0
// @description The Awayfromus web API
// @termsOfService https://awayfromus.dev/terms

// @contact.name Awayfromus Support
// @contact.url https://discord.com/invite/spqAG3Ktkv
// @contact.email thanathip.suw@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := gin.Default()

	util.LoadConfig()
	time.RunUpdateResourceExpiredDate()

	database.ConnectDB()

	app.Use(cors.Default())

	v1 := app.Group("/api/v1")

	routes.ResourceRoute(v1)
	routes.HomepageRoute(v1)
	routes.AuthRoute(v1)
	routes.DownloadRoute(v1)
	routes.UserRoute(v1)
	routes.OrderRoute(v1)
	routes.BenefitRoute(v1)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app.Run(":" + util.AppConfig.Port)
}
