package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type SignOutBtn struct {
	app.Compo
}

func (n *SignOutBtn) Render() app.UI {
	return app.Form().
		Class("d-flex").
		Attr("action", "/signout").
		Attr("method", "post").
		Body(
			app.Button().
				Class("btn btn-secondary").
				Attr("type", "submit").
				Text("Logout"),
		)
}
