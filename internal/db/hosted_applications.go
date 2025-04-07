package db

import (
	"context"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func HostedApplications() (d []*models.HostedApplication, err error) {
	var fn dbValFunc = func() any { return &models.HostedApplication{} }

	r, err := readAllCollection(fn, collectionHostedApps)

	d = make([]*models.HostedApplication, len(r))
	for i, v := range r {
		d[i] = v.(*models.HostedApplication)
	}

	return
}
func HostedApplicationByID(id string) (s *models.HostedApplication, err error) {
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
		Collection(collectionHostedApps).
		FindOne(context.TODO(), f)

	s = &models.HostedApplication{}

	err = res.Decode(s)

	return
}
func UpdateHostedApplication(h *models.HostedApplication) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"_id", h.ID}}
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collectionHostedApps).
		ReplaceOne(ctx, f, h)

	return
}
func HostedApplicationByServiceToken(serviceToken string) (s *models.HostedApplication, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"service_token", serviceToken}}
	res := c.Database(clientInstance.DbName).
		Collection(collectionHostedApps).
		FindOne(context.TODO(), f)

	s = &models.HostedApplication{}

	err = res.Decode(s)

	return
}
