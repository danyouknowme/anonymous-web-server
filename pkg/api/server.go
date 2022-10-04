package api

import (
	"github.com/danyouknowme/awayfromus/pkg/model"
)

func errorResponse(err error) model.ErrorResponse {
	error := model.ErrorResponse{
		Message: err.Error(),
	}
	return error
}
