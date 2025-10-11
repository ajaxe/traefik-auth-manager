package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type dbValFunc func() any

func readAllCollectionWithClient(c *mongo.Client, v dbValFunc, collection string) (d []any, err error) {
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
func readAllCollection(v dbValFunc, collection string) (d []any, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	return readAllCollectionWithClient(c, v, collection)
}

func deleteByIDWithClient(c *mongo.Client, id bson.ObjectID, collection string) (err error) {
	f := bson.D{{"_id", id}}

	res, err := c.Database(clientInstance.DbName).
		Collection(collection).
		DeleteMany(context.TODO(), f)

	_ = res.DeletedCount

	return
}
func deleteByID(id bson.ObjectID, collection string) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	return deleteByIDWithClient(c, id, collection)
}

func insertRecordWithClient(c *mongo.Client, u any, collection string) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collection).
		InsertOne(ctx, u)

	return
}
func insertRecord(u any, collection string) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	return insertRecordWithClient(c, u, collection)
}
