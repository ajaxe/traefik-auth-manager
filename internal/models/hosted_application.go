package models

import "go.mongodb.org/mongo-driver/v2/bson"

type HostedApplication struct {
	ID           bson.ObjectID `bson:"_id" json:"id"`
	Name         string        `json:"name"`
	ServiceToken string        `bson:"service_token" json:"serviceToken"`
	ServiceURL   string        `bson:"service_url" json:"serviceUrl"`
	Active       bool          `json:"active"`
}
type HostedAppListResult struct {
	ApiResult
	Data []*HostedApplication `json:"data"`
}