package frontend

import (
	"fmt"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
)

func HostedAppList(u string) (d models.HostedAppListResult, err error) {
	d = models.HostedAppListResult{}
	err = httpGet(buildApiURL(u, "/hosted-apps"), &d)

	if err != nil {
		return
	}

	return
}
func PostHostedApp(u string, payload, response interface{}) error {
	return httpPost(buildApiURL(u, "/hosted-apps"), payload, response)
}
func PutHostedApp(u, id string, payload, response interface{}) error {
	return httpPut(buildApiURL(u, fmt.Sprintf("/hosted-apps/%s", id)), payload, response)
}
