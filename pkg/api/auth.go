package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/models"
	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/danyouknowme/awayfromus/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

type CreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Username  string `json:"username" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req CreateUserRequest
		var user models.User
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "This username already use!"})
			return
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		newUser := models.User{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Phone:      req.Phone,
			Username:   req.Username,
			Password:   hashedPassword,
			License:    utils.GenerateLicense(req.Username),
			Resources:  []models.UserResource{},
			LastReset:  time.Now().Format(time.RFC3339),
			ResetTime:  5,
			SecretCode: utils.GenerateSecretCode(),
		}
		_, err = userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusCreated, newUser)
	}
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	FirstName  string                `json:"firstName"`
	LastName   string                `json:"lastName"`
	Email      string                `json:"email"`
	Phone      string                `json:"phone"`
	Username   string                `json:"username"`
	IsAdmin    bool                  `json:"isAdmin"`
	License    string                `json:"license"`
	Resources  []models.UserResource `json:"resources"`
	LastReset  string                `json:"lastReset"`
	ResetTime  int64                 `json:"resetTime"`
	SecretCode []string              `json:"secretCode"`
}

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func newUserResponse(user models.User) UserResponse {
	return UserResponse{
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

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req LoginUserRequest
		var user models.User
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": "Not found user with username:" + req.Username})
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = utils.CheckPassword(req.Password, user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken, err := token.CreateToken(req.Username, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response := LoginUserResponse{
			AccessToken: accessToken,
			User:        newUserResponse(user),
		}
		c.JSON(http.StatusOK, response)
	}
}
