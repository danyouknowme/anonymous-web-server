package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/models"
	"github.com/danyouknowme/awayfromus/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req CreateUserRequest
		var user models.User
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		if user.Username == req.Username {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "This username already use!"})
			return
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
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
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, newUser)
	}
}
