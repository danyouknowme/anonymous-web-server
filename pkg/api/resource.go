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

// GetResourceByName godoc
// @summary Get resource by name
// @description Get resource information by resoure name
// @tags resource
// @id GetResourceByName
// @produce json
// @param resourceName path string true "Resource name that need to get information"
// @response 200 {object} model.GetResourceByNameResponse "OK"
// @response 404 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources/{resourceName} [get]
func GetResourceByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var resource model.GetResourceByNameResponse
		resourceName := c.Param("resourceName")
		defer cancel()

		err := resourceCollection.FindOne(ctx, bson.M{"name": resourceName}).Decode(&resource)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found resource: " + resourceName)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusOK, resource)
	}
}

// GetResourceByNameAndRequiredAdmin godoc
// @summary Get resource by name required admin
// @description Get resource information by resoure name and require admin
// @tags resource
// @security ApiKeyAuth
// @id GetResourceByNameAndRequiredAdmin
// @produce json
// @param resourceName path string true "Resource name that need to get information"
// @response 200 {object} model.Resource "OK"
// @response 404 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources/admin/{resourceName} [get]
func GetResourceByNameAndRequiredAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var resource model.Resource
		resourceName := c.Param("resourceName")
		defer cancel()

		err := resourceCollection.FindOne(ctx, bson.M{"name": resourceName}).Decode(&resource)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found resource: " + resourceName)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusOK, resource)
	}
}

// AddResourceToUser godoc
// @summary Add resource to user
// @description Add resource to user and required admin
// @tags resource
// @security ApiKeyAuth
// @id AddResourceToUser
// @accept json
// @produce json
// @param Resource body model.AddResourceToUserRequest true "Username and resource information"
// @response 200 {array} model.UserResource "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources/user [post]
func AddResourceToUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.AddResourceToUserRequest
		var resource model.Resource
		var user model.User
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with username: " + req.Username)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		userResources := user.Resources

		err = resourceCollection.FindOne(ctx, bson.M{"name": req.Resource.Name}).Decode(&resource)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found resource: " + req.Resource.Name)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newResource := model.UserResource{
			Name:    req.Resource.Name,
			DayLeft: req.Resource.DayLeft,
			Status:  nil,
		}

		userResources = append(userResources, newResource)
		result := userCollection.FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"resources": userResources}})
		if result.Err() != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(result.Err()))
			return
		}

		c.JSON(http.StatusOK, userResources)
	}
}
