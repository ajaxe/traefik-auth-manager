package db

import (
	"context"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AppUsers() (d []*models.AppUser, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	cur, err := c.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		Find(ctx, bson.D{})

	if err != nil {
		return
	}
	defer cur.Close(ctx)

	d = make([]*models.AppUser, 0)
	for cur.Next(ctx) {
		var r models.AppUser
		if err = cur.Decode(&r); err != nil {
			return
		}
		d = append(d, &r)
	}

	return
}
