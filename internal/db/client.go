package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var client *mongo.Client

func NewClient() (*mongo.Client, error) {
	cfg := helpers.MustLoadDefaultAppConfig()
	return NewClientWithConfig(cfg)
}

func NewClientWithConfig(c helpers.AppConfig) (*mongo.Client, error) {
	if client != nil {
		return client, nil
	}

	sAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(c.Database.ConnectionURI).
		SetServerAPIOptions(sAPI)

	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	if err = Ping(); err != nil {
		return nil, fmt.Errorf("ping faild to return in 2sec timeout: %v", err)
	}

	return client, nil
}

func Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("ping failed with 2sec timeout: %v", err)
	}
	return nil
}

func Terminate(ctx context.Context) error {
	e := client.Disconnect(ctx)
	client = nil
	return e
}
