package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo
	users          []*models.AppUser
	allApps        []*models.HostedApplication
	appMapInternal map[string]*models.HostedApplication
}

func (u *UserList) OnNav(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""

	ctx.Async(func() {
		d, _ := frontend.UserList(b.String())
		h, _ := frontend.HostedAppList(b.String())

		ctx.Dispatch(func(c app.Context) {
			u.users = d.Data
			u.allApps = h.Data
		})
	})
}
func (u *UserList) OnMount(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""
	ctx.Handle(actionUserListReload, func(ctx app.Context, a app.Action) {
		d, _ := frontend.UserList(b.String())
		h, _ := frontend.HostedAppList(b.String())

		u.users = d.Data
		u.allApps = h.Data
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
			user:    r,
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

	for _, k := range u.allApps {
		m[strings.ToLower(k.Name)] = k
	}

	u.appMapInternal = m

	return m
}

type UserListItem struct {
	app.Compo
	user    *models.AppUser
	allApps map[string]*models.HostedApplication
}

func (ul *UserListItem) Render() app.UI {
	i := fmt.Sprintf("c%v", time.Now().UnixMilli())
	return &CardListItem{
		title: ul.user.UserName,
		actionItems: func() []app.UI {
			return []app.UI{
				&UserListItemEdit{user: ul.user},
				&UserDeleteBtn{user: ul.user},
			}
		},
		content: func() []app.UI {
			return []app.UI{
				&UserAppAssignment{
					userApps: ul.user.Applications,
					ID:       i,
					allApps:  ul.allApps,
					userId:   ul.user.ID.Hex(),
				},
			}
		},
	}
}

type UserListItemEdit struct {
	app.Compo
	user *models.AppUser
}

func (u *UserListItemEdit) Render() app.UI {
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
