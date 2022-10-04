package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/danyouknowme/awayfromus/pkg/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

// CreateUser godoc
// @summary Create user
// @description Create user
// @tags user
// @id CreateUser
// @accept json
// @produce json
// @param User body model.CreateUserRequest true "User data to be created"
// @response 200 {object} model.User "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/auth/register [post]
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.CreateUserRequest
		var user model.User
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		if user.Username == req.Username {
			err := errors.New("this username already use")
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		newUser := model.User{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Phone:      req.Phone,
			Username:   req.Username,
			Password:   hashedPassword,
			License:    util.GenerateLicense(req.Username),
			Resources:  []model.UserResource{},
			LastReset:  time.Now().Format(time.RFC3339),
			ResetTime:  5,
			SecretCode: util.GenerateSecretCode(),
		}
		_, err = userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusCreated, newUser)
	}
}

// LoginUser godoc
// @summary Login user
// @description Login user
// @tags user
// @id LoginUser
// @accept json
// @produce json
// @param User body model.LoginUserRequest true "User data to be logged in"
// @response 200 {object} model.LoginUserResponse "OK"
// @response 400 {object} model.ErrorResponse "Unauthorized"
// @response 401 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/auth/login [post]
func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.LoginUserRequest
		var user model.User
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with username" + req.Username)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = util.CheckPassword(req.Password, user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken, err := token.CreateToken(req.Username, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response := model.LoginUserResponse{
			AccessToken: accessToken,
			User:        newUserResponse(user),
		}
		c.JSON(http.StatusOK, response)
	}
}

func newUserResponse(user model.User) model.UserResponse {
	return model.UserResponse{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Phone:      user.Phone,
		Username:   user.Username,
		IsAdmin:    user.IsAdmin,
		License:    user.License,
		Resources:  user.Resources,
		LastReset:  user.LastReset,
		ResetTime:  user.ResetTime,
		SecretCode: user.SecretCode,
	}
}
