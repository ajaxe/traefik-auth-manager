package pages

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HomePage struct {
	app.Compo
	authenticated bool
}

func (h *HomePage) OnMount(ctx app.Context) {
	ctx.ObserveState(frontend.StateKeyIsAuth, &h.authenticated)
}

func (h *HomePage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{h.home()},
	}
}

func (h *HomePage) home() app.UI {
	return app.Div().
		Class("").
		Body(
			app.Div().Class("d-flex justify-content-center").
				Body(
					app.Div().Body(
						app.H1().Class("display-1 text-center").Text("Proxy Auth Manager"),
						app.P().Class("lead p-3").
							Text("Manage applications & users that require forward authentication on local Traefik instance."),
					),
				),
			app.Div().Class("d-flex justify-content-center").
				Body(
					app.If(h.authenticated, func() app.UI {
						return app.Div().Body(
							app.A().Class("btn btn-link").
								Text("Manage users").
								Href("/users"),
							app.A().Class("btn btn-link").
								Text("Manage applications").
								Href("/apps"),
						)
					},
					).Else(func() app.UI {
						return app.Div().Body(
							app.Form().
								Attr("action", "/login").
								Attr("method", "post").
								Body(
									app.Button().
										Class("btn btn-primary").
										Attr("type", "submit").
										Text("Login"),
								),
						)
					}),
				),
		)
}
