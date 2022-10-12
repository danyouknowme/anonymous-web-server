package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Accept  = "accept"
	Decline = "decline"
	Pending = "pending"
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
			resource, err := GetResourceByLabelHelper(ctx, rs.Resource)
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
			Status:           Pending,
		}

		_, err = orderCollection.InsertOne(ctx, newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusCreated, newOrder)
	}
}

// ConfirmOrder godoc
// @summary Confirm order
// @description Confirm order require admin
// @tags orders
// @security ApiKeyAuth
// @id ConfirmOrder
// @accept json
// @produce json
// @param ConfirmOrderRequest body model.ConfirmOrderRequest true "Confirm order request to be updated"
// @response 200 {object} model.MessageResponse "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/orders/confirmation [post]
func ConfirmOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var req model.ConfirmOrderRequest
		var order model.Order
		var user model.User
		var benefits []model.Benefit
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"billno": req.BillNumber}).Decode(&order)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found order with bill number: " + req.BillNumber)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if order.Status == Accept {
			err = fmt.Errorf("bill number: %s is already accepted", req.BillNumber)
			c.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"username": order.Username}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = errors.New("not found user with username: " + order.Username)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newUserResources := user.Resources
		order.Status = req.Status
		for _, rs := range order.Resources {
			if req.Status == Accept {
				newResource := model.UserResource{
					Name:    rs.ResourceName,
					DayLeft: GeneratePlanRoutine(rs.Plan),
					Status:  nil,
				}
				newUserResources = append(newUserResources, newResource)

				newBenefits := model.Benefit{
					ResourceName: rs.ResourceName,
					Price:        rs.Price,
				}
				benefits = append(benefits, newBenefits)
			}
		}

		_, err = orderCollection.UpdateOne(ctx, bson.M{"billno": order.BillNumber}, bson.M{"$set": bson.M{"status": order.Status}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{"username": order.Username}, bson.M{"$set": bson.M{"resources": newUserResources}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = UpdatePartnerBenefitsHelper(ctx, benefits)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		c.JSON(http.StatusOK, messageResponse("successfully confirm order with bill number: "+order.BillNumber))
	}
}
