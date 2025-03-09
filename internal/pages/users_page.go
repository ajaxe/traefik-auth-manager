package pages

import (
	"github.com/ajaxe/traefik-auth-manager/internal/components"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UsersPage struct {
	app.Compo
}

func (h *UsersPage) Render() app.UI {
	return &MainLayout{
		Content: []app.UI{
			app.Div().Class("row justify-content-center").Body(
				app.Div().Class("col col-md-10 col-lg-8 col-xl-6").Body(
					components.AppUserList(),
				),
			),
			components.AppUserEditModal(),
		},
	}
}
