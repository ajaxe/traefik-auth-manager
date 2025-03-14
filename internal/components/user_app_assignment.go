package components

import (
	"sort"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserAppAssignment struct {
	app.Compo
	ID       string
	userApps []*models.ApplicationIdentifier
	allApps  map[string]*models.HostedApplication
}

func (u *UserAppAssignment) Render() app.UI {
	existing := len(u.userApps) > 0

	return app.If(!existing, func() app.UI {
		return app.Div().
			Class("collapse").
			ID(u.ID).
			Body(
				app.Text("No applications assigned to this user."),
			)
	}).Else(func() app.UI {
		return app.Div().
			ID(u.ID).
			Body(
				u.listApps()...,
			)
	})
}
func (u *UserAppAssignment) listApps() []app.UI {
	l := []app.UI{}
	a := make(map[string]bool)

	for _, r := range u.userApps {
		a[strings.ToLower(r.Name)] = true
	}

	keys := []string{}
	for k := range u.allApps {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		c := "btn-outline-secondary"
		if _, ok := a[k]; ok {
			c = "btn-primary"
		}

		l = append(l, app.Button().Class("me-1 btn btn-sm "+c).Text(u.allApps[k].Name))
	}
	return l
}
