package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type NavBar struct {
	app.Compo
	authenticated bool
}

func (n *NavBar) OnNav(ctx app.Context) {
	ctx.ObserveState(frontend.StateKeyIsAuth, &n.authenticated)
}

func (n *NavBar) Render() app.UI {
	return app.Nav().
		Class("navbar navbar-expand-lg navbar-dark bg-primary").
		Body(
			app.Div().Class("container").
				Body(
					n.brandName(),
					app.If(n.authenticated, func() app.UI {
						return n.navToggler()
					}), app.If(n.authenticated, func() app.UI {
						return n.navItems()
					}),
				),
		)
}

func (n NavBar) brandName() app.UI {
	return &AppName{}
}
func (n NavBar) navItems() app.UI {
	return app.Div().
		Class("collapse navbar-collapse").
		ID("navbarSupportedContent").
		Body(
			AppNavBarItems(NavListOptions{
				TextColor: "text-white",
				ListCSS:   "navbar-nav me-auto mb-2 mb-lg-0",
			}),
			AppSignoutBtn(),
		)
}
func (n NavBar) navToggler() app.UI {
	return app.Button().
		Class("navbar-toggler").
		Class("collapsed").
		Type("button").
		DataSet("bs-toggle", "collapse").
		DataSet("bs-target", "#navbarSupportedContent").
		Aria("controls", "navbarSupportedContent").
		Aria("expanded", "false").
		Aria("label", "Toggle navigation").
		Body(
			app.Span().Class("navbar-toggler-icon"),
		)
}
