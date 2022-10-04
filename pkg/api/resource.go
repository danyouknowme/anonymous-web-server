package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetResourcesInformation godoc
// @summary Get Resources
// @description Get all resource information
// @tags resource
// @id GetResourcesInformation
// @produce json
// @response 200 {array} model.AllResourceResponse "OK"
// @response 500 {object} model.ErrorResponse "Not Found"
// @router /api/v1/resources [get]
func GetAllResourcesInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var response []model.AllResourceResponse
		defer cancel()

		results, err := resourceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var resource model.AllResourceResponse
			if err = results.Decode(&resource); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			response = append(response, resource)
		}
		c.JSON(http.StatusOK, response)
	}
}

// CreateNewResource godoc
// @summary Create new resource
// @description Create a new resource
// @tags resource
// @security ApiKeyAuth
// @id CreateNewResource
// @accept json
// @produce json
// @param Resource body model.Resource true "Resource data to be created"
// @response 201 {object} model.Resource "Created"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources [post]
func CreateNewResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.Resource
		var ownedResource model.Resource
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		err := resourceCollection.FindOne(ctx, bson.M{"name": req.Name}).Decode(&ownedResource)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		if ownedResource.Name == req.Name {
			err := errors.New("this resource already have")
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newResource := model.Resource{
			IsPublish:   req.IsPublish,
			Name:        req.Name,
			Label:       req.Label,
			Description: req.Description,
			Document:    req.Document,
			Video:       req.Video,
			Thumbnail:   req.Thumbnail,
			Images:      req.Images,
			Plan:        req.Plan,
			PatchNotes:  req.PatchNotes,
		}

		_, err = resourceCollection.InsertOne(ctx, newResource)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusCreated, newResource)
	}
}

// UpdateResource godoc
// @summary Update resource
// @description Update resource information
// @tags resource
// @security ApiKeyAuth
// @id UpdateResource
// @accept json
// @produce json
// @param Resource body model.Resource true "Resource data to be updated"
// @response 200 {object} model.Resource "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources [patch]
func UpdateResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.Resource
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		resourceUpdated := bson.M{
			"is_publish":  req.IsPublish,
			"name":        req.Name,
			"label":       req.Label,
			"description": req.Description,
			"document":    req.Document,
			"video":       req.Video,
			"thumbnail":   req.Thumbnail,
			"images":      req.Images,
			"plan":        req.Plan,
			"patch_notes": req.PatchNotes,
		}
		result := resourceCollection.FindOneAndUpdate(ctx, bson.M{"name": req.Name}, bson.M{"$set": resourceUpdated})
		if result.Err() != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(result.Err()))
			return
		}

		c.JSON(http.StatusOK, resourceUpdated)
	}
}
