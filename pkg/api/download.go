package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetDownloadResource godoc
// @summary Get download resource
// @description Get download resource by resoure name
// @tags download
// @security ApiKeyAuth
// @id GetDownloadResource
// @produce json
// @param resourceName path string true "Resource name that need to get download"
// @response 200 {array} model.GetDownloadResourceResponse "OK"
// @response 401 {object} model.ErrorResponse "Unauthorized"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources/download/{resourceName} [get]
func GetDownloadResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var resource model.Resource
		var user model.User
		var response []model.GetDownloadResourceResponse
		resourceName := c.Param("resourceName")
		username := c.MustGet(authorizationPayloadKey).(*token.Payload).Username
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with username: " + username)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = resourceCollection.FindOne(ctx, bson.M{"name": resourceName}).Decode(&resource)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found resource " + resourceName)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		for _, rs := range user.Resources {
			if rs.Name == resource.Name {
				for _, pn := range resource.PatchNotes {
					downloadResource := model.GetDownloadResourceResponse{
						Version:  pn.Version,
						Download: pn.Download,
					}

					response = append(response, downloadResource)
				}
				c.JSON(http.StatusOK, response)
				return
			}
		}

		err = errors.New("unallowed to get download of " + resourceName)
		c.JSON(http.StatusUnauthorized, errorResponse(err))
	}
}

// DownloadResource godoc
// @summary Download resource
// @description Download resource by resoure name
// @tags download
// @security ApiKeyAuth
// @id DownloadResource
// @accept json
// @produce json
// @param Download body model.DownloadResourceRequest true "Resource name that need to create download"
// @response 200 {object} model.Download "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 401 {object} model.ErrorResponse "Unauthorized"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/resources/download/ [post]
func DownloadResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.DownloadResourceRequest
		var user model.User
		var resource model.Resource
		// var response []model.GetDownloadResourceResponse
		// resourceName := c.Param("resourceName")
		username := c.MustGet(authorizationPayloadKey).(*token.Payload).Username
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with username: " + username)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = resourceCollection.FindOne(ctx, bson.M{"name": req.ResourceName}).Decode(&resource)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found resource " + req.ResourceName)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		for _, rs := range user.Resources {
			if rs.Name == resource.Name {
				newDownload := model.Download{
					Username:     user.Username,
					ResourceName: resource.Name,
					Date:         time.Now().Format(time.RFC3339),
				}
				_, err = downloadCollection.InsertOne(ctx, newDownload)
				if err != nil {
					c.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}

				c.JSON(http.StatusOK, newDownload)
				return
			}
		}

		err = errors.New("user: " + user.Username + " unallowed to download " + resource.Name)
		c.JSON(http.StatusUnauthorized, errorResponse(err))
	}
}
