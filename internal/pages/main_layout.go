package pages

import (
	"github.com/ajaxe/traefik-auth-manager/internal/components"
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MainLayout struct {
	app.Compo
	Content []app.UI
}

func (m *MainLayout) Render() app.UI {
	return m.container()
}

func (m *MainLayout) OnNav(ctx app.Context) {
	baseURL := app.Window().URL()
	baseURL.Path = ""

	ctx.Async(func() {
		sess, err := frontend.CheckAuth(baseURL.String())

		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				frontend.AppContext(ctx).SetIsAuth(false)
				if app.Window().URL().Path != "/home" {
					ctx.Navigate("/home")
				}
			} else {
				frontend.AppContext(ctx).
					SetSession(sess).
					SetIsAuth(true)
			}
		})
	})
}

func (m *MainLayout) container() app.UI {
	return app.Main().
		Body(
			components.AppNavBar(),
			app.Div().Class("container mt-3").
				ID("main-container").
				Body(
					m.Content...,
				),
			components.AppCodeUpdate(),
		)
}
