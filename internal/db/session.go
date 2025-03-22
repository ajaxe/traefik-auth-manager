package db

import (
	"context"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertSession(s models.Session) (id string, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	res, err := c.Database(clientInstance.DbName).
		Collection(collectionSession).
		InsertOne(context.TODO(), s)
	if err != nil {
		return
	}

	id = res.InsertedID.(bson.ObjectID).Hex()
	return
}

func SessionByID(id string) (s *models.Session, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	hex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	f := bson.D{{"_id", hex}}
	res := c.Database(clientInstance.DbName).
		Collection(collectionSession).
		FindOne(context.TODO(), f)

	s = &models.Session{}

	err = res.Decode(s)

	return
}

func DeleteSessionByID(id string) (err error) {

	hex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	err = deleteByID(hex, collectionSession)
	if err != nil {
		return
	}

	return
}
