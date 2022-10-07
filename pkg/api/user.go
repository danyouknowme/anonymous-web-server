package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyUserAdmin(ctx *gin.Context, username string) (bool, error) {
	var user model.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return false, err
	}

	if user.IsAdmin {
		return true, nil
	}
	return false, nil
}

func UpdateUserResourceExpiredDate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var user model.User
		if err = results.Decode(&user); err != nil {
			return err
		}

		resources := []model.UserResource{}

		for _, rs := range user.Resources {
			resourceUpdated := model.UserResource{}
			if rs.DayLeft > 0 {
				resourceUpdated = model.UserResource{
					Name:    rs.Name,
					DayLeft: rs.DayLeft - 1,
					Status:  rs.Status,
				}
			} else {
				resourceUpdated = model.UserResource{
					Name:    rs.Name,
					DayLeft: rs.DayLeft,
					Status:  rs.Status,
				}
			}
			resources = append(resources, resourceUpdated)
		}

		result := userCollection.FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"resources": resources}})
		if result.Err() != nil {
			return result.Err()
		}
	}
	return nil
}

// CheckLicense godoc
// @summary Check license
// @description Check license and update resource status
// @tags users
// @id CheckLicense
// @accept json
// @produce json
// @param License body model.CheckLicenseRequest true "License key and resource name that need to update"
// @response 200 {array} model.UserResource "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/license [post]
func CheckLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.CheckLicenseRequest
		var user model.User
		var resource model.Resource
		ip := c.ClientIP()
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"license": req.License}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with license: " + req.License)
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

		resources := []model.UserResource{}

		for _, rs := range user.Resources {
			var status *string
			if rs.Name == resource.Name {
				status = &ip
			} else {
				status = rs.Status
			}
			userResource := model.UserResource{
				Name:    rs.Name,
				Status:  status,
				DayLeft: rs.DayLeft,
			}
			resources = append(resources, userResource)
		}

		result := userCollection.FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"resources": resources}})
		if result.Err() != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(result.Err()))
			return
		}
		c.JSON(http.StatusOK, resources)
	}
}

// ResetIP godoc
// @summary Reset ip resources
// @description Reset ip of all user resource status
// @tags users
// @security ApiKeyAuth
// @id ResetIP
// @produce json
// @response 200 {array} model.UserResource "OK"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/ip/reset [post]
func ResetIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user model.User
		var resources []model.UserResource
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

		if user.ResetTime > 0 {
			err = errors.New("wait for " + strconv.FormatInt(int64(user.ResetTime), 10) + " minutes for reset again")
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		for _, rs := range user.Resources {
			resource := model.UserResource{
				Name:    rs.Name,
				Status:  nil,
				DayLeft: rs.DayLeft,
			}
			resources = append(resources, resource)
		}

		result := userCollection.FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"resources": resources}})
		if result.Err() != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(result.Err()))
			return
		}
		c.JSON(http.StatusOK, resources)
	}
}

// GetUserData godoc
// @summary Get user data
// @description Get user data required admin
// @tags users
// @security ApiKeyAuth
// @id GetUserData
// @produce json
// @param username path string true "Username of user that need to update data"
// @response 200 {object} model.GetUserDataResponse "OK"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/{username} [get]
func GetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user model.User
		username := c.Param("username")
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

		response := model.GetUserDataResponse{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Email:      user.Email,
			Phone:      user.Phone,
			Username:   user.Username,
			License:    user.License,
			Resources:  user.Resources,
			SecretCode: user.SecretCode,
		}

		c.JSON(http.StatusOK, response)
	}
}

// UpdateUserData godoc
// @summary Update user data
// @description Update user data required admin
// @tags users
// @security ApiKeyAuth
// @id UpdateUserData
// @accept json
// @produce json
// @param User body model.UpdateUserDataRequest true "User data that will be updated"
// @response 200 {object} model.MessageResponse "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users [patch]
func UpdateUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.UpdateUserDataRequest
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

		updatedUser := model.User{
			Username:   user.Username,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Phone:      req.Phone,
			License:    req.License,
			Resources:  req.Resources,
			Password:   user.Password,
			IsAdmin:    user.IsAdmin,
			LastReset:  time.Now().Format(time.RFC3339),
			ResetTime:  user.ResetTime,
			SecretCode: user.SecretCode,
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{}, bson.M{"$set": updatedUser})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		message := fmt.Sprintf("%s successfully updated", user.Username)
		c.JSON(http.StatusOK, messageResponse(message))
	}
}
