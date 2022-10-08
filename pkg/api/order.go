package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllOrders godoc
// @summary Get all orders
// @description Get all orders require admin
// @tags orders
// @security ApiKeyAuth
// @id GetAllOrders
// @produce json
// @response 200 {array} model.Order "OK"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/orders [get]
func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orders []model.Order
		defer cancel()

		findOptions := options.Find()
		findOptions.SetSort(bson.M{"status": -1})

		results, err := orderCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var order model.Order
			if err = results.Decode(&order); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, orders)
	}
}
