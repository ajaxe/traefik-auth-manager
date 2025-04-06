package frontend

import (
	"fmt"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	StateKeyIsAuth        = "isAuth"
	StateKeyUserSession   = "user-session"
	StateKeyUserList      = "user-list"
	StateKeyHostedAppList = "hosted-app-list"
)

func NewAppContext(ctx app.Context) AppContext {
	return AppContext{
		Context: ctx,
	}
}

type AppContext struct {
	app.Context
}

func (c AppContext) Session() models.Session {
	var s models.Session
	c.GetState(StateKeyUserSession, &s)
	return s
}
func (c AppContext) SetSession(s models.Session) AppContext {
	c.SetState(StateKeyUserSession, s)
	return c
}
func (c AppContext) SetIsAuth(a bool) AppContext {
	c.SetState(StateKeyIsAuth, a)
	return c
}
func (c AppContext) IsAuth() bool {
	var a bool
	c.GetState(StateKeyIsAuth, &a)
	return a
}
func appBaseURL() string {
	b := app.Window().URL()
	b.Path = ""
	return b.String()
}
func buildApiURL(b, p string) string {
	return fmt.Sprintf("%s/api/%s", strings.TrimSuffix(b, "/"), strings.TrimPrefix(p, "/"))
}
func (c AppContext) LoadData(key string) {
	switch key {
	case StateKeyUserList:
		c.loadUserList()
	case StateKeyHostedAppList:
		c.loadHostedAppList()
	default: // do nothing
	}
}
func (c AppContext) loadHostedAppList() {
	b := appBaseURL()
	c.Async(func() {
		l, _ := HostedAppList(b)
		c.SetState(StateKeyHostedAppList, l.Data)
	})
}
func (c AppContext) loadUserList() {
	b := appBaseURL()

	c.Async(func() {
		d, _ := UserList(b)
		h, _ := HostedAppList(b)
		c.SetState(StateKeyUserList, NewUserListViewData(d.Data, h.Data))
	})
}
func (c AppContext) UpdateHostedApp(id string, payload models.HostedApplication) (err error) {
	u := appBaseURL()

	r := &models.ApiResult{}
	err = PutHostedApp(u, id, payload, &r)
	return
}
func (c AppContext) ToggleUserApp(userId, appID string, selected bool, cb func()) {
	b := appBaseURL()
	c.Async(func() {
		var err error
		r := models.ApiResult{}
		if selected {
			err = RemoveUserApp(userId, appID, b, r)
		} else {
			err = AssignUserApp(userId, appID, b, r)
		}
		if err != nil {
			helpers.AppLogf("%v", err)
		} else if cb != nil {
			cb()
		}
	})
}
