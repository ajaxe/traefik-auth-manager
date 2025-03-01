package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type LoginBtn struct {
	app.Compo
}

func (l *LoginBtn) Render() app.UI {
	return app.Form().
		Attr("action", "/login").
		Attr("method", "post").
		Body(
			app.Button().
				Class("btn btn-primary").
				Attr("type", "submit").
				Text("Login"),
		)
}
