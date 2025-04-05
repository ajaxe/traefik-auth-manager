package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type HostedAppAddBtn struct {
	app.Compo
}
func (u *HostedAppAddBtn) Render() app.UI {
	return app.Div().
		Class("d-flex mb-3").
		Body(
			app.Button().
				Class("btn btn-primary btn-sm ms-auto").
				Text("Add Application").
				OnClick(func(ctx app.Context, e app.Event) {
					ctx.NewAction(actionHostedAppAdd)
				}),
		)
}