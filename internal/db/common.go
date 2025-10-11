package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type dbValFunc func() any

var tracer = otel.Tracer("db")

func readAllCollectionWithClient(c *mongo.Client, v dbValFunc, collection string, ct ...context.Context) (d []any, err error) {
	var ctx context.Context
	if len(ct) > 0 {
		ctx = ct[0]
	} else {
		ctx = context.Background()
	}

	ctx, span := tracer.Start(ctx, "readAllCollection")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBCollectionNameKey.String(collection),
		semconv.DBOperationName("find"),
	)

	cur, err := c.Database(clientInstance.DbName).
		Collection(collection).
		Find(ctx, bson.D{})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		r := v()
		if err = cur.Decode(r); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		d = append(d, r)
	}
	span.SetAttributes(attribute.Int("db.rows_affected", len(d)))
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
