package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type dbValFunc func() any

func packageTracer(ct context.Context) trace.Tracer {
	return trace.SpanFromContext(ct).TracerProvider().Tracer("db")
}

func readCollectionWithClient(c *mongo.Client, v dbValFunc, collection string, ct ...context.Context) (d []any, err error) {
	var ctx context.Context
	if len(ct) > 0 {
		ctx = ct[0]
	} else {
		ctx = context.TODO()
	}

	ctx, span := packageTracer(ctx).Start(ctx, "db.readCollection")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBCollectionNameKey.String(collection),
		semconv.DBOperationName("find"),
	)

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

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
	span.SetAttributes(attribute.Int("db.readCollection.rows_affected", len(d)))
	return
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

func insertRecordWithClient(c *mongo.Client, u any, collection string, ct ...context.Context) (err error) {
	var ctx context.Context
	if len(ct) > 0 {
		ctx = ct[0]
	} else {
		ctx = context.TODO()
	}

	ctx, span := packageTracer(ctx).Start(ctx, "db.insertRecord")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBCollectionNameKey.String(collection),
		semconv.DBOperationName("insertOne"),
	)

	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collection).
		InsertOne(ctx, u)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return
}
func insertRecord(u any, collection string) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	return insertRecordWithClient(c, u, collection)
}
