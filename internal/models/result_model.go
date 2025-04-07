package models

import (
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApiResult struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}
type ApiIDResult struct {
	ApiResult
	ID bson.ObjectID `json:"id"`
}

func NewApiIDResult(id bson.ObjectID) *ApiIDResult {
	return &ApiIDResult{
		ApiResult: ApiResult{
			Success: true,
		},
		ID: id,
	}
}
func NewGenericErrApiResult(err error) ApiResult {
	helpers.AppLogf("error: %s", err.Error())
	return ApiResult{
		Success:      false,
		ErrorMessage: "Something went wrong.",
	}
}
