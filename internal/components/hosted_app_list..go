package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
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
	ctx.Handle(actionHostedAppReload, h.reloadData)
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
			title:       i.Name,
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
			helpers.AppLogf("name=%s, url=%s", i.happ.Name, i.happ.ServiceURL)
		},
	}
	return func() []app.UI {
		return []app.UI{b}
	}
}
func (h *HostedAppList) reloadData(ctx app.Context, a app.Action) {
	helpers.AppLogf("reloadData: hosted-app list")
	h.loadDataInternal(ctx)
}
func (h *HostedAppList) loadDataInternal(ctx app.Context) {
	b := app.Window().URL()
	b.Path = ""
	ctx.Async(func() {
		l, _ := frontend.HostedAppList(b.String())

		ctx.Dispatch(func(c app.Context) {
			h.apps = l.Data
			helpers.AppLogf("loadDataInternal: hosted-app list")
		})
	})
}

type HostedAppListItem struct {
	app.Compo
	happ         *models.HostedApplication
	originalData *models.HostedApplication
	ReadOnly     bool
	Compact      bool
	formResult   *models.ApiResult
}

func (h *HostedAppListItem) OnMount(ctx app.Context) {
	ctx.Handle(actionHostedAppEdit, h.handleOnEdit)
}

func (h *HostedAppListItem) Render() app.UI {
	helpers.AppLogf("happ: Name=%s:Compact=%v:ReadOnly=%v", h.happ.Name, h.Compact, h.ReadOnly)
	return app.Form().Class("row").
		Body(
			h.serviceTokenUI(),
			h.activeCheckbox(),
			h.serviceURLUI(),
			h.editActions(),
		)
}
func (h *HostedAppListItem) serviceTokenUI() app.UI {
	id := "ha-svc-token"
	return app.Div().Class("col-md-9").
		Body(
			&FormControl{
				Compact: h.Compact,
				Content: func() []app.UI {
					return []app.UI{
						&FormText{
							ID:          id,
							Value:       h.happ.ServiceToken,
							BindTo:      &h.happ.ServiceToken,
							InputType:   "text",
							ReadOnly:    h.ReadOnly,
							Placeholder: "Not set",
						},
						&FormLabel{
							For:   id,
							Label: "Service Token",
						},
					}
				},
			},
		)
}
func (h *HostedAppListItem) activeCheckbox() app.UI {
	return app.Div().Class("col-md-3").
		Body(
			&FormCheckbox{
				label:    "Active",
				BindTo:   h.happ.Active,
				Value:    h.happ.Active,
				role:     "switch",
				Disabled: h.ReadOnly,
			},
		)
}

func (h *HostedAppListItem) serviceURLUI() app.UI {
	id := "ha-svc-url"
	return app.Div().Class("col-md-12").
		Body(
			&FormControl{
				Compact: h.Compact,
				Content: func() []app.UI {
					return []app.UI{
						&FormText{
							ID:        id,
							Value:     h.happ.ServiceURL,
							BindTo:    &h.happ.ServiceURL,
							InputType: "text",
							ReadOnly:  h.ReadOnly,
						},
						&FormLabel{
							For:   id,
							Label: "Service URL",
						},
					}
				},
			},
		)
}
func (h *HostedAppListItem) handleOnEdit(ctx app.Context, a app.Action) {
	id, ok := a.Value.(string)
	if !ok || id != h.happ.ID.Hex() {
		h.cancel()
		return
	}
	o := *h.happ
	h.originalData = &o

	helpers.AppLogf("Are happ equal orig == copy? = %v", h.happ == h.originalData)

	h.readonlyView(false)
}
func (h *HostedAppListItem) editActions() app.UI {
	return app.If(h.ReadOnly, func() app.UI { return app.Div() }).
		Else(func() app.UI {
			return app.Div().Class("col-md-12").
				Body(
					app.Button().Class("btn btn-primary").Text("Save").
						OnClick(h.onSave),
					app.Button().Class("btn btn-link ms-1").Text("Cancel").
						OnClick(h.onCancel),
					h.formResultMessage(),
				)
		})
}

func (h *HostedAppListItem) formResultMessage() app.UI {
	if h.formResult == nil {
		return app.Span()
	}
	ico := "bi bi-check-circle text-success"
	co := "text-success"
	m := "Saved"

	if !h.formResult.Success {
		ico = "bi bi-x-circle text-danger"
		co = "text-danger"
		m = h.formResult.ErrorMessage
	}

	return app.Span().Class("ms-2 fw-bold").Body(
		app.I().Class(ico),
		app.Span().Class("ms-1 "+co).Text(m),
	)
}
func (h *HostedAppListItem) readonlyView(v bool) {
	h.ReadOnly = v
	h.Compact = v
}
func (h *HostedAppListItem) cancel() {
	if h.originalData != nil {
		h.happ = h.originalData
		h.originalData = nil
	}
	h.readonlyView(true)
}
func (h *HostedAppListItem) onCancel(ctx app.Context, e app.Event) {
	h.cancel()
	e.PreventDefault()
	ctx.Update()
}
func (h *HostedAppListItem) onSave(ctx app.Context, e app.Event) {
	e.PreventDefault()
	u := app.Window().URL()
	u.Path = ""

	helpers.AppLogf("OnSave: original: %v", h.originalData)
	helpers.AppLogf("OnSave: current: %v", h.happ)
	r := &models.ApiResult{}
	err := frontend.PutHostedApp(u.String(), h.originalData.ID.Hex(), h.happ, &r)

	if err != nil {
		helpers.AppLogf("h-apps onSave error: %v", err)
		return
	}
	h.cancel()
	ctx.NewAction(actionHostedAppReload)
}
