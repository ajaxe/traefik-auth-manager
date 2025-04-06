package frontend

import (
	"fmt"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	StateKeyIsAuth      = "isAuth"
	StateKeyUserSession = "user-session"
	StateKeyUserList    = "user-list"
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
func buildApiURL(b, p string) string {
	return fmt.Sprintf("%s/api/%s", strings.TrimSuffix(b, "/"), strings.TrimPrefix(p, "/"))
}

func (c AppContext) LoadHostedAppList() {
	c.Async(func() {
		b := app.Window().URL()
		b.Path = ""
		l, _ := HostedAppList(b.String())
		c.SetState(StateKeyUserList, l.Data)
	})
}
func (c AppContext) UpdateHostedApp(id string, payload models.HostedApplication) (err error) {
	u := app.Window().URL()
	u.Path = ""

	r := &models.ApiResult{}
	err = PutHostedApp(u.String(), id, payload, &r)
	return
}
