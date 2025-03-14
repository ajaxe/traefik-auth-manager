package frontend

import "github.com/ajaxe/traefik-auth-manager/internal/models"

func HostedAppList(u string) (d models.HostedAppListResult, err error) {
	d = models.HostedAppListResult{}
	err = httpGet(buildApiURL(u, "/hosted-apps"), &d)

	if err != nil {
		return
	}

	return
}
