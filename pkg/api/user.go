package api

import (
	"context"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
			}
			resources = append(resources, resourceUpdated)
		}

		result := userCollection.FindOneAndUpdate(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"resources": resources}})
		if result.Err() != nil {
			return err
		}
	}
	return nil
}
