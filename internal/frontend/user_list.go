package frontend

import (
	"github.com/ajaxe/traefik-auth-manager/internal/models"
)

func UserList(u string) (d []models.AppUser, err error) {
	d = []models.AppUser{}
	err = httpGet(buildApiURL(u, "/app-users"), &d)

	if err != nil {
		return
	}

	return
}
