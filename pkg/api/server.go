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

func errorResponse(err error) model.ErrorResponse {
	error := model.ErrorResponse{
		Message: err.Error(),
	}
	return error
}
