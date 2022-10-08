package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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

func UpdateUserResetTimeOnMinute(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"reset_time": user.ResetTime}})
	if err != nil {
		return err
	}

	return nil
}

func GetResourceByLabelHelper(ctx context.Context, resourceLabel string) (resource model.Resource, err error) {
	err = resourceCollection.FindOne(ctx, bson.M{"label": resourceLabel}).Decode(&resource)
	if err != nil {
		return
	}
	return
}

func FindPlanHelper(requestPlan string, resource model.Resource) (plan model.Plan, err error) {
	for _, p := range resource.Plan {
		if p.Name == requestPlan {
			return p, nil
		}
	}
	err = fmt.Errorf("plan %s not found in resource %s", requestPlan, resource.Label)
	return plan, err
}

func GenerateBillNumberHelper(ctx context.Context) (billNumber string, err error) {
	m := 6
	c := "0"

	count, err := orderCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return "", err
	}
	orderCount := strconv.FormatInt(count+1, 10)

	if n := utf8.RuneCountInString(orderCount); n < m {
		orderCount = strings.Repeat(c, m-n) + orderCount
	}

	return orderCount, nil
}

func GeneratePlanRoutine(plan string) int64 {
	switch plan {
	case "MONTH":
		return 30
	case "YEAR":
		return 365
	default:
		return -1
	}
}
