package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type UserAddBtn struct {
	app.Compo
}

func (u *UserAddBtn) Render() app.UI {
	return app.Div().
		Class("d-flex mb-3").
		Body(
			app.Button().
				Class("btn btn-primary btn-sm ms-auto").
				Text("Add User").
				OnClick(func(ctx app.Context, e app.Event) {
					ctx.NewAction(actionUserAdd)
				}),
		)
}
