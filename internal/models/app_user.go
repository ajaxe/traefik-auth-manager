package models

import "go.mongodb.org/mongo-driver/v2/bson"

type ApplicationIdentifier struct {
	HostAppId bson.ObjectID `json:"hostedAppId"`
	Name      string        `json:"name"`
}

type AppUser struct {
	ID           bson.ObjectID            `bson:"_id" json:"id"`
	UserName     string                   `json:"userName"`
	Password     string                   `json:"-"`
	Active       bool                     `json:"active"`
	Applications []*ApplicationIdentifier `json:"applications"`
}

type AppUserListResult struct {
	ApiResult
	Data []*AppUser `json:"data"`
}

type AppUserChange struct {
	AppUser
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
