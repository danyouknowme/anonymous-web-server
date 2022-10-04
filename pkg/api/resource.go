package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var resourceCollection *mongo.Collection = database.GetCollection(database.DB, "resources")

type AllResourceResponse struct {
	Name      string       `json:"name"`
	Label     string       `json:"label"`
	Thumbnail string       `json:"thumbnail"`
	Plan      []model.Plan `json:"plan"`
	IsPublish bool         `json:"is_publish"`
}

// GetResourcesInformation godoc
// @summary Get Resources
// @description Get all resource information
// @tags resouces
// @id GetResourcesInformation
// @produce json
// @response 200 {array} model.Resource "OK"
// @response 500 {object} model.ErrorResponse "Not Found"
// @router /api/v1/resources [get]
func GetAllResourcesInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var response []AllResourceResponse
		defer cancel()

		results, err := resourceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var resource AllResourceResponse
			if err = results.Decode(&resource); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			response = append(response, resource)
		}
		c.JSON(http.StatusOK, response)
	}
}
