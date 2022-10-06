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
