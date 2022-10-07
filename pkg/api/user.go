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
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

		if user.ResetTime != 5 {
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

		_, err = userCollection.UpdateOne(ctx, bson.M{"username": user.Username}, bson.M{"$set": updatedUser})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		message := fmt.Sprintf("%s successfully updated", user.Username)
		c.JSON(http.StatusOK, messageResponse(message))
	}
}

// RemoveUserResource godoc
// @summary Remove user resource
// @description Remove user resource by index
// @tags users
// @security ApiKeyAuth
// @id RemoveUserResource
// @produce json
// @param resourceIndex path int true "Index of resource that need to remove"
// @response 200 {array} model.UserResource "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/resource/{resourceIndex} [delete]
func RemoveUserResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user model.User
		var userResource []model.UserResource
		username := c.MustGet(authorizationPayloadKey).(*token.Payload).Username
		defer cancel()

		resourceIndex, err := strconv.Atoi(c.Param("resourceIndex"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with username: " + username)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		for i, rs := range user.Resources {
			if i != resourceIndex {
				userResource = append(userResource, rs)
			} else {
				if rs.DayLeft > 0 {
					err = fmt.Errorf("%s still have %d days left", rs.Name, rs.DayLeft)
					c.JSON(http.StatusBadRequest, errorResponse(err))
					return
				}
			}
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"resources": userResource}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		c.JSON(http.StatusOK, userResource)
	}
}

// GetAllUsername godoc
// @summary Get all username
// @description Get all username require admin
// @tags users
// @security ApiKeyAuth
// @id GetAllUsername
// @produce json
// @response 200 {array} string "OK"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/usernames [get]
func GetAllUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var usernames []string
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var user model.User
			if err = results.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			usernames = append(usernames, user.Username)
		}

		c.JSON(http.StatusOK, usernames)
	}
}

// GetUserResetTime godoc
// @summary Get user reset time
// @description Get user reset time
// @tags users
// @security ApiKeyAuth
// @id GetUserResetTime
// @produce json
// @response 200 {object} model.GetUserResetTimeResponse "OK"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/reset-time [get]
func GetUserResetTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user model.User
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

		response := model.GetUserResetTimeResponse{
			Timer:    user.ResetTime,
			CanReset: user.ResetTime == 5,
		}

		c.JSON(http.StatusOK, response)
	}
}

// UpdateUserResetTime godoc
// @summary Update user reset time
// @description Update user reset time
// @tags users
// @security ApiKeyAuth
// @id UpdateUserResetTime
// @produce json
// @response 200 {object} model.MessageResponse "OK"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/users/reset-time [patch]
func UpdateUserResetTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user model.User
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

		if user.ResetTime != 5 {
			err = fmt.Errorf("timer is still running wait for %d minute(s)", user.ResetTime)
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Minutes().Do(func() {
			if user.ResetTime != 0 {
				user.ResetTime--
				err = UpdateUserResetTimeOnMinute(user)
				if err != nil {
					c.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
			} else {
				user.ResetTime = 5
				err = UpdateUserResetTimeOnMinute(user)
				if err != nil {
					c.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
				s.Stop()
			}
		})
		s.StartAsync()

		c.JSON(http.StatusOK, messageResponse("ok"))
	}
}
