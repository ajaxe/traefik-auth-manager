package models

import "go.mongodb.org/mongo-driver/v2/bson"

type ApplicationIdentifier struct {
	HostAppId bson.ObjectID
}
type AppUser struct {
	UserName     string
	Password     string
	Active       bool
	Applications []ApplicationIdentifier
}
