package main

// "log"

import (
	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/routes"
	"github.com/danyouknowme/awayfromus/pkg/utils"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	utils.LoadConfig()

	database.ConnectDB()

	routes.ResourceRoute(app)
	routes.HomepageRoute(app)

	app.Run(":" + utils.AppConfig.Port)
}
