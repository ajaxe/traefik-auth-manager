package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	collectionAppUser    = "app_users"
	collectionSession    = "session"
	collectionHostedApps = "hosted_applications"
)
const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
)

type AppClient struct {
	Client *mongo.Client
	DbName string
}

var clientInstance AppClient

func NewClient() (*mongo.Client, error) {
	cfg := helpers.MustLoadDefaultAppConfig()
	return newClientWithConfig(cfg)
}
func newAppClient() AppClient {
	if clientInstance.Client == nil {
		cfg := helpers.MustLoadDefaultAppConfig()
		newClientWithConfig(cfg)
	}
	return clientInstance
}

func newClientWithConfig(c helpers.AppConfig) (*mongo.Client, error) {
	if clientInstance.Client != nil {
		return clientInstance.Client, nil
	}

	sAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(c.Database.ConnectionURI).
		SetServerAPIOptions(sAPI)

	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	if err = pingClient(client, context.Background()); err != nil {
		return nil, fmt.Errorf("ping faild to return in 2sec timeout: %v", err)
	}

	clientInstance.Client = client
	clientInstance.DbName = c.Database.DbName

	return client, nil
}

func Ping(ctx context.Context) error {
	return pingClient(clientInstance.Client, ctx)
}
func pingClient(c *mongo.Client, ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "pingClient")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	span.SetAttributes(
		semconv.DBSystemMongoDB,
		semconv.DBNamespaceKey.String(clientInstance.DbName),
		semconv.DBOperationName("ping"),
	)

	if c == nil {
		err := fmt.Errorf("client must be instantiated before calling Ping")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err := c.Ping(ctx, readpref.Primary())
	if err != nil {
		e := fmt.Errorf("ping failed with 2sec timeout: %v", err)
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}
	return nil
}

func Terminate(ctx context.Context) error {
	client := clientInstance.Client
	if client == nil {
		log.Print("db client not instantiated, nothing to disconnect")
		return nil
	}
	e := client.Disconnect(ctx)
	client = nil
	return e
}
