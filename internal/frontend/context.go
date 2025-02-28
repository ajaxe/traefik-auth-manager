package frontend

import (
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	StateKeyIsAuth     = "isAuth"
	StatKeyUserSession = "user-session"
)

type AppContext app.Context

func (c AppContext) Session() models.Session {
	var s models.Session
	app.Context(c).GetState(StatKeyUserSession, &s)
	return s
}
func (c AppContext) SetSession(s models.Session) AppContext {
	app.Context(c).SetState(StatKeyUserSession, s)
	return c
}
func (c AppContext) SetIsAuth(a bool) AppContext {
	app.Context(c).SetState(StateKeyIsAuth, a)
	return c
}
func (c AppContext) IsAuth() bool {
	var a bool
	app.Context(c).GetState(StateKeyIsAuth, &a)
	return a
}
