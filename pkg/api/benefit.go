package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPartnerBenefits godoc
// @summary Get partner benefits
// @description Get partner benefits
// @tags benefits
// @security ApiKeyAuth
// @id GetPartnerBenefits
// @produce json
// @response 200 {array} model.Benefit "OK"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/benefits [get]
func GetPartnerBenefits() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var response []model.Benefit
		defer cancel()

		results, err := benefitCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var benefit model.Benefit
			if err = results.Decode(&benefit); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			response = append(response, benefit)
		}

		fmt.Println(response)

		c.JSON(http.StatusOK, response)
	}
}

// ClearPartnerBenefit godoc
// @summary Clear partner benefit
// @description Clear partner benefit by resoure name
// @tags benefits
// @security ApiKeyAuth
// @id ClearPartnerBenefit
// @accept json
// @produce json
// @param Download body model.ClearPartnerBenefitRequest true "Resource name that need to update partner benefit"
// @response 200 {object} model.MessageResponse "OK"
// @response 400 {object} model.ErrorResponse "Bad Request"
// @response 404 {object} model.ErrorResponse "Not Found"
// @response 500 {object} model.ErrorResponse "Internal Server Error"
// @router /api/v1/benefits/clear [post]
func ClearPartnerBenefit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var benefit model.Benefit
		var req model.ClearPartnerBenefitRequest
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if validationErr := validate.Struct(&req); validationErr != nil {
			c.JSON(http.StatusBadRequest, errorResponse(validationErr))
			return
		}

		err := benefitCollection.FindOne(ctx, bson.M{"resource_name": req.ResourceName}).Decode(&benefit)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
		}

		_, err = benefitCollection.UpdateOne(ctx, bson.M{"resource_name": req.ResourceName}, bson.M{"$set": bson.M{"price": 0}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response := fmt.Sprintf("successfully clear benefits of %s", req.ResourceName)
		c.JSON(http.StatusOK, messageResponse(response))
	}
}
