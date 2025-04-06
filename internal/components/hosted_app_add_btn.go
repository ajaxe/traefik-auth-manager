package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HostedAppAddBtn struct {
	app.Compo
	showForm bool
}

func (u *HostedAppAddBtn) Render() app.UI {
	return app.Div().Body(
		app.Div().
			Class("d-flex mb-3").
			Body(
				app.Button().
					Class("btn btn-primary btn-sm ms-auto").
					Text("Add Application").
					OnClick(func(ctx app.Context, e app.Event) {
						u.showForm = true
					}),
			),
		u.newApp(),
	)
}
func (u *HostedAppAddBtn) newApp() app.UI {
	if !u.showForm {
		return app.Div()
	}
	return newHostedAppCardItem(models.HostedApplication{}, HostedAppCardOptions{
		ReadOnly: false,
		Compact:  false,
	})
}
