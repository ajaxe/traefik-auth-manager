package db

import (
	"github.com/ajaxe/traefik-auth-manager/internal/models"
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
