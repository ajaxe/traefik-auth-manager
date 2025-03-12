package components

import (
	"fmt"
	"time"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo
	users []*models.AppUser
}

func (u *UserList) OnNav(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""

	ctx.Async(func() {
		d, _ := frontend.UserList(b.String())

		ctx.Dispatch(func(c app.Context) {
			u.users = d.Data

			ctx.Update()
		})
	})
}

func (u *UserList) Render() app.UI {
	return app.Div().Body(
		u.userListItems()...,
	)
}

func (u *UserList) userListItems() []app.UI {
	l := []app.UI{}
	for _, r := range u.users {
		l = append(l, &UserListItem{
			user: r,
		})
	}
	return l
}

type UserListItem struct {
	app.Compo
	user *models.AppUser
}

func (ul *UserListItem) Render() app.UI {
	i := fmt.Sprintf("c%v", time.Now().UnixMilli())
	return app.Div().Class("card").Style("margin-bottom", "10px").
		Body(
			app.Div().Class("card-body").
				Body(
					app.Div().Class("card-title d-flex").
						Body(
							app.Div().Class("me-auto").
								Style("cursor", "pointer", "padding-top", "5px").
								DataSet("bs-toggle", "collapse").
								DataSet("bs-target", "#"+i).
								Role("button").
								Aria("expanded", "false").
								Aria("controls", i).
								Body(
									app.Span().Class("h5").Text(ul.user.UserName),
									app.I().Class("bi bi-arrow-right ms-2"),
								),
							&UserListItemEdit{user: ul.user},
						),
					app.Div().Class("collapse").ID(i).
						Text("more info"),
				),
		)
}

type UserListItemEdit struct {
	app.Compo
	user *models.AppUser
}

func (u *UserListItemEdit) Render() app.UI {
	return app.Div().
		Body(
			app.Button().Class("btn btn-light").
				DataSet("user-id", u.user.ID.String()).
				OnClick(func(ctx app.Context, e app.Event){
					ctx.NewActionWithValue(actionUserEdit, u.user)
				}).
				Body(
					app.I().Class("bi bi-pencil-square"),
				),
		)
}
