package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type LoginAvatar struct {
	app.Compo
	displayCss string
	s          models.Session
}

func (a *LoginAvatar) OnMount(ctx app.Context) {
	ctx.ObserveState(frontend.StateKeyUserSession, &a.s)
}
func (a *LoginAvatar) Render() app.UI {
	return app.Div().Class("align-self-center " + a.displayCss).
		Body(
			app.Div().Class("d-inline-block shadow-lg dropdown rounded-pill").
				Body(
					app.Img().Width(48).Class("img-thumbnail rounded-pill").
						Title(a.s.User.Name).
						Src(a.s.User.Picture),
				),
		)
}
