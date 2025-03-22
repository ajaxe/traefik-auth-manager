package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HostedAppList struct {
	app.Compo
	apps []*models.HostedApplication
}

func (u *HostedAppList) OnNav(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""

	ctx.Async(func() {
		h, _ := frontend.HostedAppList(b.String())

		ctx.Dispatch(func(c app.Context) {
			u.apps = h.Data
		})
	})
}
func (h *HostedAppList) Render() app.UI {
	return app.Div().Body(
		h.listItems()...,
	)
}
func (h *HostedAppList) listItems() []app.UI {
	items := []app.UI{}

	for _, i := range h.apps {
		items = append(items, &CardListItem{
			title: i.Name,
		})
	}
	return items
}
