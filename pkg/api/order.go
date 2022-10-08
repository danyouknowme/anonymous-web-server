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

// AddOrder godoc
// @summary Add order
// @description Add order by requested resource and plan
// @tags orders
// @security ApiKeyAuth
// @id AddOrder
// @accept json
// @produce json
// @param OrderRequest body model.AddOrderRequest true "Order request to be created"
// @response 201 {object} model.Order "Created"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/orders [post]
func AddOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.AddOrderRequest
		var totalPrice float64 = 0
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		orderResources := []model.OrderResource{}
		for _, rs := range req.RequestOrder {
			resource, err := GetResourceByNameHelper(ctx, rs.Resource)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			plan, err := FindPlanHelper(rs.Plan, resource)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalPrice += plan.Price

			newOrderResource := model.OrderResource{
				ResourceName:  resource.Name,
				ResourceLabel: resource.Label,
				Plan:          plan.Name,
				Price:         plan.Price,
			}

			orderResources = append(orderResources, newOrderResource)
		}

		billNumber, err := GenerateBillNumberHelper(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newOrder := model.Order{
			BillNumber:       billNumber,
			Username:         req.Username,
			Resources:        orderResources,
			TransactionImage: req.TransactionImage,
			Date:             time.Now().Format(time.RFC3339),
			TotalPrice:       totalPrice,
			Status:           "pending",
		}

		_, err = orderCollection.InsertOne(ctx, newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusCreated, newOrder)
	}
}
