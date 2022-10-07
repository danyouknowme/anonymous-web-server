package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetHomepageInformation godoc
// @summary Get Homepage
// @description Get homepage information
// @tags homepage
// @id GetHomepageInformation
// @produce json
// @response 200 {object} model.Homepage "OK"
// @response 404 {object} model.ErrorResponse "Not Found"
// @router /api/v1/homepage [get]
func GetHomepageInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var homepageInfo model.Homepage
		defer cancel()

		err := homepageCollection.FindOne(ctx, bson.M{}).Decode(&homepageInfo)
		if err != nil {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusOK, homepageInfo)
	}
}

// UpdateHomepageInformation godoc
// @summary Update Homepage
// @description Update homepage information
// @tags homepage
// @security ApiKeyAuth
// @id UpdateHomepageInformation
// @accept json
// @produce json
// @param Homepage body model.Homepage true "Homepage data to be updated"
// @response 200 {object} model.Homepage "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/homepage [patch]
func UpdateHomepageInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.Homepage
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

// GetCounterState godoc
// @summary Get counter state
// @description Get counter state information to show on homepage
// @tags homepage
// @id GetCounterState
// @produce json
// @response 200 {object} model.GetCounterStateResponse "OK"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/homepage/counter [get]
func GetCounterState() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userCount, err := userCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		downloadCount, err := downloadCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		orderCount, err := orderCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response := model.GetCounterStateResponse{
			Users:     userCount,
			Downloads: downloadCount,
			Orders:    orderCount,
		}

		c.JSON(http.StatusOK, response)
	}
}
