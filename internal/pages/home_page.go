package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HomePage struct {
	app.Compo
}

func (h *HomePage) Render() app.UI {
	return &MainLayout{
		Content: app.Div().Text("home"),
	}
}
