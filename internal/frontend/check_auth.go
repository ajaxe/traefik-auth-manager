package frontend

import (
	"github.com/ajaxe/traefik-auth-manager/internal/models"
)

func CheckAuth(u string) (models.Session, error) {
	var s models.Session
	err := httpGet(buildApiURL(u, "/check"), &s)

	if err != nil {
		return models.Session{}, err
	}

	return s, nil
}
