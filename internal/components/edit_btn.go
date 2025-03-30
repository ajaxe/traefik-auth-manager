package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type EditBtn struct {
	app.Compo
	onClick func(app.Context, app.Event)
}

func (b *EditBtn) Render() app.UI {
	e := app.Button().Class("btn btn-light").
		Body(
			app.I().Class("bi bi-pencil-square"),
		)

	if b.onClick != nil {
		e.OnClick(func(ctx app.Context, e app.Event) {
			b.onClick(ctx, e)
		})
	}

	return e
}
