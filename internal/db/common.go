package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type dbValFunc func() any

func readAllCollection(v dbValFunc, collection string) (d []any, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	cur, err := c.Database(clientInstance.DbName).
		Collection(collection).
		Find(ctx, bson.D{})

	if err != nil {
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		r := v()
		if err = cur.Decode(r); err != nil {
			return
		}
		d = append(d, r)
	}

	return
}
