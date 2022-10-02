package main

// "log"

import (
	"log"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/utils"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	database.ConnectDB(config.MONGO_URI)

	app.Run(":" + config.PORT)
}
