package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CopyValueBtn struct {
	app.Compo
	Value   string
	Visible bool
}

func (c *CopyValueBtn) Render() app.UI {
	btn := app.Button().Class("btn btn-sm btn-light").
		Body(
			app.I().Class("bi bi-copy"),
		).
		OnClick(func(ctx app.Context, e app.Event) {
			frontend.NewAppContext(ctx).CopyToClipboard(c.Value)
			e.PreventDefault()
		})
	return app.If(c.Visible && len(c.Value) > 0, func() app.UI {
		return btn
	}).Else(func() app.UI {
		return app.Span()
	})
}
