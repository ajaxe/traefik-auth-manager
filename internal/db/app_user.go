package db

import (
	"context"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AppUsers() (d []*models.AppUser, err error) {
	var fn dbValFunc = func() any { return &models.AppUser{} }

	r, err := readAllCollection(fn, collectionAppUser)

	d = make([]*models.AppUser, len(r))
	for i, v := range r {
		d[i] = v.(*models.AppUser)
	}

	return
}

func UpdatePassword(u *models.AppUser) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"_id", u.ID}}
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"Password", u.Password}}}}

	_, err = c.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		UpdateOne(ctx, f, update)

	return
}
