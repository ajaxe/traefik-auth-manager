package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type DeleteBtn struct {
	app.Compo
	ID       string
	confirm  bool
	onDelete func(ctx app.Context)
}

func (u *DeleteBtn) Render() app.UI {
	return app.If(!u.confirm, func() app.UI {
		return u.deleteBtn()
	}).Else(func() app.UI {
		return u.deleteConfirmation()
	})
}

func (u *DeleteBtn) deleteBtn() app.UI {
	return app.Button().
		ID(u.ID).
		Class("btn btn-light").
		OnClick(func(ctx app.Context, e app.Event) {
			u.confirm = true
		}).
		Body(
			app.I().Class("bi bi-trash3"),
		)
}

func (u *DeleteBtn) deleteConfirmation() app.UI {
	return app.Div().
		Class("d-flex align-items-center").
		Body(
			app.Span().Text("Confirm?"),
			app.Span().
				Body(
					app.Button().
						Class("btn btn-sm btn-success me-1 ms-1").
						Body(
							app.I().Class("bi bi-trash3"),
						).
						OnClick(func(ctx app.Context, e app.Event) {
							u.confirm = false
							u.onDelete(ctx)
						}),
					app.Button().
						Class("btn btn-sm btn-danger").
						Body(
							app.I().Class("bi bi-x"),
						).
						OnClick(func(ctx app.Context, e app.Event) {
							u.confirm = false
						}),
				),
		)
}
