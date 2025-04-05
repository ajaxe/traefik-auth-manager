package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HostedAppList struct {
	app.Compo
	apps []*models.HostedApplication
}

func (h *HostedAppList) OnNav(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""

	h.loadDataInternal(ctx)
}
func (h *HostedAppList) OnMount(ctx app.Context) {
	ctx.Handle(actionHostedAppReload, func(ctx app.Context, a app.Action) {
		b := app.Window().URL()
		b.Path = ""
		l, _ := frontend.HostedAppList(b.String())

		h.apps = l.Data
		ctx.Update()
	})
}
func (h *HostedAppList) Render() app.UI {
	return app.Div().Body(
		h.listItems()...,
	)
}
func (h *HostedAppList) listItems() []app.UI {
	items := []app.UI{}

	for _, i := range h.apps {
		itm := &HostedAppListItem{
			happ:     i,
			ReadOnly: true,
			Compact:  true,
		}
		items = append(items, &CardListItem{
			Title:       i.Name,
			actionItems: h.itemActions(itm),
			content: func() []app.UI {
				return []app.UI{itm}
			},
		})
	}
	return items
}
func (h *HostedAppList) itemActions(i *HostedAppListItem) func() []app.UI {
	b := &EditBtn{
		onClick: func(ctx app.Context, e app.Event) {
			ctx.NewActionWithValue(actionHostedAppEdit, i.happ.ID.Hex())
		},
	}
	return func() []app.UI {
		return []app.UI{b}
	}
}

func (h *HostedAppList) loadDataInternal(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""
	ctx.Async(func() {
		l, _ := frontend.HostedAppList(b.String())

		ctx.Dispatch(func(ctx app.Context) {
			h.apps = l.Data
			ctx.Update()
		})
	})
}
