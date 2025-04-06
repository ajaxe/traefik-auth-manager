package components

import (
	"fmt"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserDeleteBtn struct {
	app.Compo
	user    *frontend.AppUserView
	id      string
	confirm bool
}

func (u *UserDeleteBtn) Render() app.UI {
	u.id = fmt.Sprintf("u-del-btn-%s", u.user.ID.Hex())
	return app.If(!u.confirm, func() app.UI {
		return u.deleteBtn()
	}).Else(func() app.UI {
		return u.deleteConfirmation()
	})
}

func (u *UserDeleteBtn) deleteBtn() app.UI {
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
}

func (u *UserDeleteBtn) deleteConfirmation() app.UI {
	b := app.Window().URL()
	b.Path = ""
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
							r := models.ApiResult{}
							_ = frontend.RemoveUser(u.user.ID.Hex(), b.String(), &r)
							frontend.NewAppContext(ctx).LoadData(frontend.StateKeyUserList)
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
