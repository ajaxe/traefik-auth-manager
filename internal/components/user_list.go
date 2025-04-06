package components

import (
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo
	data           frontend.UserListViewData
	appMapInternal map[string]*models.HostedApplication
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
			user:    &frontend.AppUserView{AppUser: *r},
			allApps: u.appMap(),
		})
	}
	return l
}
func (u *UserList) appMap() map[string]*models.HostedApplication {
	if u.appMapInternal != nil {
		return u.appMapInternal
	}
	m := make(map[string]*models.HostedApplication)

	for _, k := range u.data.Apps {
		m[strings.ToLower(k.Name)] = k
	}

	u.appMapInternal = m

	return m
}
