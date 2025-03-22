package models

import "go.mongodb.org/mongo-driver/v2/bson"

type ApplicationIdentifier struct {
	HostAppId bson.ObjectID `bson:"host_app_id" json:"hostAppId"`
	Name      string        `json:"name"`
}

type AppUser struct {
	ID           bson.ObjectID            `bson:"_id" json:"id"`
	UserName     string                   `bson:"username" json:"userName"`
	Password     string                   `bson:"password" json:"-"`
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
