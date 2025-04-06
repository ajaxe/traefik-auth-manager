package frontend

import "github.com/ajaxe/traefik-auth-manager/internal/models"

type AppUserView struct {
	models.AppUser
}
type ApplicationIdentifierView struct {
	models.ApplicationIdentifier
}

type UserListViewData struct {
	Users []*models.AppUser
	Apps  []*models.HostedApplication
}
