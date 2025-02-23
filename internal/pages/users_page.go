package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UsersPage struct {
	app.Compo
}

func (h *UsersPage) Render() app.UI {
	return &MainLayout{
		Content: app.Div().Text("Users"),
	}
}
