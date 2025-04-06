package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo
	data frontend.UserListViewData
}

func (u *UserList) OnMount(ctx app.Context) {
	ctx.ObserveState(frontend.StateKeyUserList, &u.data)
}

func (u *UserList) Render() app.UI {
	return app.Div().Body(
		u.userListItems()...,
	)
}

func (u *UserList) userListItems() []app.UI {
	l := []app.UI{}
	for _, r := range u.data.Users {
		l = append(l, &UserListItem{
			user: r,
		})
	}
	return l
}
