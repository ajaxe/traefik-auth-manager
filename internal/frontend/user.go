package frontend

import (
	"github.com/ajaxe/traefik-auth-manager/internal/models"
)

func UserList(u string) (d models.AppUserListResult, err error) {
	d = models.AppUserListResult{}
	err = httpGet(buildApiURL(u, "/app-users"), &d)

	if err != nil {
		return
	}

	return
}

func PostUser(u string, payload, response interface{}) error {
	return httpPost(buildApiURL(u, "/app-users"), payload, response)
}
func PutUser(id, u string, payload, response interface{}) error {
	return httpPut(buildApiURL(u, "/app-users/"+id), payload, response)
}
