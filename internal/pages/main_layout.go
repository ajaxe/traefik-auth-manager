package pages

import (
	"github.com/ajaxe/traefik-auth-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MainLayout struct {
	app.Compo
	Content app.UI
}

func (m *MainLayout) Render() app.UI {
	return m.container()
}

func (m *MainLayout) container() app.UI {
	return app.Main().
		Body(
			components.AppNavBar(),
			app.Div().Class("container").
				ID("main-container").
				Body(
					m.Content,
				),
			components.AppCodeUpdate(),
		)
}
