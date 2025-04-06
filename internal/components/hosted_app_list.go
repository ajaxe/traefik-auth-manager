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

func (h *HostedAppList) OnMount(ctx app.Context) {
	ctx.ObserveState(frontend.StateKeyHostedAppList, &h.apps)
}
func (h *HostedAppList) Render() app.UI {
	return app.Div().Body(
		h.listItems()...,
	)
}
func (h *HostedAppList) listItems() []app.UI {
	items := []app.UI{}

	for _, i := range h.apps {
		items = append(items, newHostedAppCardItem(*i))
	}
	return items
}

func newHostedAppCardItem(h models.HostedApplication, o ...HostedAppCardOptions) *CardListItem {
	if len(o) == 0 {
		o = append(o, HostedAppCardOptions{
			ReadOnly: true,
			Compact:  true,
		})
	}
	itm := &HostedAppListItem{
		Happ:                 &h,
		HostedAppCardOptions: o[0],
	}
	var a []app.UI
	if h.Name != "" {
		a = itemActions(itm)
	}
	title := h.Name
	if title == "" {
		title = "New Application"
	}
	return &CardListItem{
		Title:       h.Name,
		actionItems: a,
		content:     []app.UI{itm},
	}
}
func itemActions(i *HostedAppListItem) []app.UI {
	b := &EditBtn{
		onClick: func(ctx app.Context, e app.Event) {
			ctx.NewActionWithValue(actionHostedAppEdit, i.Happ.ID.Hex())
		},
	}
	return []app.UI{b}
}
