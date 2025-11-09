package db

import (
	"context"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewHostedAppDataAccess() func(context.Context) *HostedAppDataAccess {
	return func(ctx context.Context) *HostedAppDataAccess {
		c := newAppClient()
		return &HostedAppDataAccess{
			client: c.Client,
			ctx:    ctx,
		}
	}
}

type HostedAppDataAccess struct {
	client *mongo.Client
	ctx    context.Context
}

func (h *HostedAppDataAccess) HostedApplications() (d []*models.HostedApplication, err error) {
	var fn dbValFunc = func() any { return &models.HostedApplication{} }

	r, err := readCollectionWithClient(h.client, fn, collectionHostedApps, h.ctx)

	d = make([]*models.HostedApplication, len(r))
	for i, v := range r {
		d[i] = v.(*models.HostedApplication)
	}

	return
}
func (h *HostedAppDataAccess) HostedApplicationByID(id string) (s *models.HostedApplication, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	hex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	ctx, span := packageTracer(h.ctx).Start(h.ctx, "db.by_id")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBCollectionNameKey.String(collectionHostedApps),
		semconv.DBOperationName("findOne"),
	)

	f := bson.D{{"_id", hex}}
	res := c.Database(clientInstance.DbName).
		Collection(collectionHostedApps).
		FindOne(ctx, f)

	s = &models.HostedApplication{}

	err = res.Decode(s)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return
}
func (h *HostedAppDataAccess) UpdateHostedApplication(d *models.HostedApplication) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"_id", d.ID}}

	ctx, span := packageTracer(h.ctx).Start(h.ctx, "db.update")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBCollectionNameKey.String(collectionHostedApps),
		semconv.DBOperationName("replaceOne"),
	)

	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collectionHostedApps).
		ReplaceOne(ctx, f, d)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return
}
func (h *HostedAppDataAccess) HostedApplicationByServiceToken(serviceToken string) (s *models.HostedApplication, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"service_token", serviceToken}}

	ctx, span := packageTracer(h.ctx).Start(h.ctx, "db.by_service_token")
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBCollectionNameKey.String(collectionHostedApps),
		semconv.DBOperationName("findOne"),
	)

	res := c.Database(clientInstance.DbName).
		Collection(collectionHostedApps).
		FindOne(ctx, f)

	s = &models.HostedApplication{}

	err = res.Decode(s)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return
}
func (h *HostedAppDataAccess) InsertHostedApplication(u *models.HostedApplication) (id bson.ObjectID, err error) {

	id = bson.NewObjectID()
	u.ID = id

	err = insertRecordWithClient(h.client, u, collectionHostedApps, h.ctx)
	
	return
}

func DeleteHostedAppByID(id bson.ObjectID) error {
	return deleteByID(id, collectionHostedApps)
}
