package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HostedAppAddBtn struct {
	app.Compo
	showForm bool
	instance *models.HostedApplication
}

func (u *HostedAppAddBtn) Render() app.UI {
	helpers.AppLogf("HostedAppAddBtn render")
	return app.Div().Body(
		app.Div().
			Class("d-flex mb-3").
			Body(
				app.Button().
					Class("btn btn-primary btn-sm ms-auto").
					Text("Add Application").
					OnClick(func(ctx app.Context, e app.Event) {
						u.showForm = true
						e.PreventDefault()
					}),
			),
		u.newApp(),
	)
}
func (u *HostedAppAddBtn) hide() {
	u.showForm = false
	u.instance = nil
}
func (u *HostedAppAddBtn) newApp() app.UI {
	if !u.showForm {
		return app.Div()
	}
	if u.instance == nil {
		u.instance = &models.HostedApplication{}
	}
	helpers.AppLogf("newApp instance[%p]: %v", u.instance, u.instance)
	var itm = HostedAppListItem{
		Happ:   u.instance,
		onSave: u.onSave,
		onCancel: func(ctx app.Context) {
			u.hide()
		},
		HostedAppCardOptions: HostedAppCardOptions{
			ReadOnly: false,
			Compact:  false,
		},
	}
	return newHostedAppCardItemWithItem(&itm)
}
func (u *HostedAppAddBtn) onSave(ctx app.Context) {
	helpers.AppLogf("onSave instance[%p]: %v", u.instance, u.instance)
	err := frontend.NewAppContext(ctx).AddHostedApp(*u.instance)
	defer func() { u.hide() }()
	if err != nil {
		return
	}

	frontend.NewAppContext(ctx).LoadData(frontend.StateKeyHostedAppList)
}
