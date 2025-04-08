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
func newHostedAppCardItem(h models.HostedApplication) *CardListItem {
	nested := &HostedAppListItem{
		Happ: &h,
		HostedAppCardOptions: HostedAppCardOptions{
			ReadOnly: true,
			Compact:  true,
		},
	}
	return newHostedAppCardItemWithItem(nested)
}
func newHostedAppCardItemWithItem(nested *HostedAppListItem) *CardListItem {
	var a []app.UI
	if nested.Happ.Name != "" {
		a = itemActions(nested)
	}
	title := nested.Happ.Name
	if title == "" {
		title = "New Application"
	}
	return &CardListItem{
		Title:       nested.Happ.Name,
		actionItems: a,
		content:     []app.UI{nested},
	}
}
func itemActions(i *HostedAppListItem) []app.UI {
	b := &EditBtn{
		onClick: func(ctx app.Context, e app.Event) {
			ctx.NewActionWithValue(actionHostedAppEdit, i.Happ.ID.Hex())
		},
	}
	d := &HostedAppDeleteBtn{
		Happ: i.Happ,
	}
	return []app.UI{b, d}
}
