package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var homepageCollection *mongo.Collection = database.GetCollection(database.DB, "homepage")

func GetHomepageInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var homepageInfo models.Homepage
		defer cancel()

		err := homepageCollection.FindOne(ctx, bson.M{}).Decode(&homepageInfo)
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		fmt.Println(homepageInfo.ResourceName)

		c.JSON(http.StatusOK, homepageInfo)
	}
}
