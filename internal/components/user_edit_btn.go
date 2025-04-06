package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserEditBtn struct {
	app.Compo
	user *frontend.AppUserView
}

func (u *UserEditBtn) Render() app.UI {
	return app.Div().
		Class("me-1").
		Body(
			app.Button().Class("btn btn-light").
				OnClick(func(ctx app.Context, e app.Event) {
					ctx.NewActionWithValue(actionUserEdit, u.user)
				}).
				Body(
					app.I().Class("bi bi-pencil-square"),
				),
		)
}
