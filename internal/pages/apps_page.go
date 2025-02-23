package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type AppsPage struct {
	app.Compo
}

func (h *AppsPage) Render() app.UI {
	return &MainLayout{
		Content: app.Div().Text("Applications"),
	}
}
