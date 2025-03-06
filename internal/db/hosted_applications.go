package db

import "github.com/ajaxe/traefik-auth-manager/internal/models"

func HostedApplications() (d []*models.HostedApplication, err error) {
	var fn dbValFunc = func() any { return &models.HostedApplication{} }

	r, err := readAllCollection(fn, collectionHostedApps)

	d = make([]*models.HostedApplication, len(r))
	for i, v := range r {
		d[i] = v.(*models.HostedApplication)
	}

	return
}
