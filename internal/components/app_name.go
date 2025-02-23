package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type AppName struct {
	app.Compo
}

func (l *AppName) Render() app.UI {
	return app.A().
		Class("navbar-brand").
		Href("#").
		Text("Proxy Auth Manager")
}
