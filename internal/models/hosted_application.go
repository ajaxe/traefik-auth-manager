package models

import "go.mongodb.org/mongo-driver/v2/bson"

type HostedApplication struct {
	ID           bson.ObjectID `bson:"_id" json:"id"`
	Name         string        `json:"name"`
	ServiceToken string        `json:"serviceToken"`
	ServiceURL   string        `json:"serviceUrl"`
	Active       bool          `json:"active"`
}
