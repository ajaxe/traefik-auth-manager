package db

import (
	"context"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewAppUserDataAccess() func(context.Context) *AppUserDataAccess {
	return func(ctx context.Context) *AppUserDataAccess {
		c := newAppClient()
		return &AppUserDataAccess{
			client: c.Client,
			ctx:    ctx,
		}
	}
}

type AppUserDataAccess struct {
	client *mongo.Client
	ctx    context.Context
}

func (c *AppUserDataAccess) AppUsers() (d []*models.AppUser, err error) {
	var fn dbValFunc = func() any { return &models.AppUser{} }

	r, err := readAllCollectionWithClient(c.client, fn, collectionAppUser, c.ctx)

	d = make([]*models.AppUser, len(r))
	for i, v := range r {
		d[i] = v.(*models.AppUser)
	}

	return
}
func (c *AppUserDataAccess) AppUserByID(id string) (s *models.AppUser, err error) {
	hex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	f := bson.D{{"_id", hex}}
	res := c.client.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		FindOne(context.TODO(), f)

	s = &models.AppUser{}

	err = res.Decode(s)

	return
}
func (c *AppUserDataAccess) UpdatePassword(u *models.AppUser) (err error) {
	f := bson.D{{"_id", u.ID}}
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"password", u.Password}}}}

	_, err = c.client.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		UpdateOne(ctx, f, update)

	return
}
func (c *AppUserDataAccess) UpdateUserHostedApps(u *models.AppUser) (err error) {
	f := bson.D{{"_id", u.ID}}
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"applications", u.Applications}}}}

	_, err = c.client.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		UpdateOne(ctx, f, update)

	return
}
func (c *AppUserDataAccess) InsertAppUser(u *models.AppUser) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	u.ID = id

	err = insertRecordWithClient(c.client, u, collectionAppUser)
	return
}
func (c *AppUserDataAccess) DeleteAppUserByID(id bson.ObjectID) error {
	return deleteByIDWithClient(c.client, id, collectionAppUser)
}
func (c *AppUserDataAccess) AppUserByUsername(username string) (s *models.AppUser, err error) {
	f := bson.D{{"username", username}}
	res := c.client.Database(clientInstance.DbName).
		Collection(collectionAppUser).
		FindOne(context.TODO(), f)

	s = &models.AppUser{}

	err = res.Decode(s)

	return
}
