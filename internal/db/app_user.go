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
func AppUserByID(id string) (s *models.AppUser, err error) {
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
		Collection(collectionAppUser).
		FindOne(context.TODO(), f)

	s = &models.AppUser{}

	err = res.Decode(s)

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
func UpdateUserHostedApps(u *models.AppUser) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"_id", u.ID}}
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"Applications", u.Applications}}}}

	_, err = c.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		UpdateOne(ctx, f, update)

	return
}

func InsertAppUser(u *models.AppUser) (id bson.ObjectID, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	id = bson.NewObjectID()
	u.ID = id

	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		InsertOne(ctx, u)

	return
}
func DeleteAppUserByID(id bson.ObjectID) error {
	return deleteByID(id, collectionAppUser)
}
