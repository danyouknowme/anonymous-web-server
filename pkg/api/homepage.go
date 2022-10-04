package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var homepageCollection *mongo.Collection = database.GetCollection(database.DB, "homepage")
var validate = validator.New()

func GetHomepageInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var homepageInfo models.Homepage
		defer cancel()

		err := homepageCollection.FindOne(ctx, bson.M{}).Decode(&homepageInfo)
		if err != nil {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		fmt.Println(homepageInfo.ResourceName)

		c.JSON(http.StatusOK, homepageInfo)
	}
}

func UpdateHomepageInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req models.Homepage
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		homepageUpdated := bson.M{
			"title":         req.Title,
			"resourceName":  req.ResourceName,
			"resourceLabel": req.ResourceLabel,
			"description":   req.Description,
		}
		_, err := homepageCollection.UpdateOne(ctx, bson.M{}, bson.M{"$set": homepageUpdated})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusOK, homepageUpdated)
	}
}
