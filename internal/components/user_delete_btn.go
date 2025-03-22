package components

import (
	"fmt"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserDeleteBtn struct {
	app.Compo
	user    *models.AppUser
	id      string
	confirm bool
}

func (u *UserDeleteBtn) OnMount(ctx app.Context) {
	helpers.AppLogf("btn id: %s", u.id)
}
func (u *UserDeleteBtn) Render() app.UI {
	u.id = fmt.Sprintf("u-del-btn-%s", u.user.ID.Hex())
	return app.If(!u.confirm, func() app.UI {
		return app.Button().
			ID(u.id).
			Class("btn btn-light").
			OnClick(func(ctx app.Context, e app.Event) {
				helpers.AppLog("clicked delete")
				u.confirm = true
			}).
			Body(
				app.I().Class("bi bi-trash3"),
			)
	}).Else(func() app.UI {
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
							),
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
	})
}
