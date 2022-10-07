package api

import (
	"github.com/danyouknowme/awayfromus/pkg/database"
	"github.com/danyouknowme/awayfromus/pkg/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

var homepageCollection *mongo.Collection = database.GetCollection(database.DB, "homepage")
var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")
var resourceCollection *mongo.Collection = database.GetCollection(database.DB, "resources")
var downloadCollection *mongo.Collection = database.GetCollection(database.DB, "downloads")
var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

func messageResponse(message string) model.MessageResponse {
	response := model.MessageResponse{
		Message: message,
	}
	return response
}

func errorResponse(err error) model.ErrorResponse {
	error := model.ErrorResponse{
		Error: err.Error(),
	}
	return error
}
