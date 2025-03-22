package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
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
			content: func() []app.UI {
				return []app.UI{
					&HostedAppListItem{
						happ: i,
					},
				}
			},
		})
	}
	return items
}

type HostedAppListItem struct {
	app.Compo
	happ *models.HostedApplication
}

func (h *HostedAppListItem) Render() app.UI {
	helpers.AppLogf("happ: %s:%s:%s", h.happ.ServiceToken, h.happ.ServiceURL, h.happ.Name)
	return app.Form().Class("row").
		Body(
			h.serviceTokenUI(),
			h.activeCheckbox(),
			h.serviceURLUI(),
		)
}
func (h *HostedAppListItem) serviceTokenUI() app.UI {
	id := "ha-svc-token"
	return app.Div().Class("col-md-6").
		Body(
			&FormControl{
				Content: func() []app.UI {
					return []app.UI{
						&FormText{
							ID:        id,
							Value:     h.happ.ServiceToken,
							BindTo:    &h.happ.ServiceToken,
							InputType: "text",
						},
						&FormLabel{
							For:   id,
							Label: "Service Token",
						},
					}
				},
			},
		)
}
func (h *HostedAppListItem) activeCheckbox() app.UI {
	return app.Div().Class("col-md-6").
		Body(
			&FormCheckbox{
				label:  "Active",
				bindTo: h.happ.Active,
				value:  h.happ.Active,
				role:   "switch",
			},
		)
}

func (h *HostedAppListItem) serviceURLUI() app.UI {
	id := "ha-svc-url"
	return app.Div().Class("col-md-12").
		Body(
			&FormControl{
				Content: func() []app.UI {
					return []app.UI{
						&FormText{
							ID:        id,
							Value:     h.happ.ServiceURL,
							BindTo:    &h.happ.ServiceURL,
							InputType: "text",
						},
						&FormLabel{
							For:   id,
							Label: "Service Token",
						},
					}
				},
			},
		)
}
