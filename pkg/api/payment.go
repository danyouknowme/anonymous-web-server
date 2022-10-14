package api

import (
	"context"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetPaymentHistory godoc
// @summary Get payment history
// @description Get payment history of user
// @tags payments
// @security ApiKeyAuth
// @id GetPaymentHistory
// @produce json
// @response 200 {array} model.Order "OK"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/payments [get]
func GetPaymentHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orders []model.Order
		username := c.MustGet(authorizationPayloadKey).(*token.Payload).Username
		defer cancel()

		results, err := orderCollection.Find(ctx, bson.M{"username": username})
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
