package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

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
	f := []app.UI{
		h.serviceNameUI(),
		h.serviceTokenUI(),
		h.activeCheckbox(),
		h.serviceURLUI(),
		h.editActions(),
	}

	if h.ReadOnly {
		f = f[1:]
	}

	return app.Form().Class("row").
		Body(f...)
}

func (h *HostedAppListItem) serviceNameUI() app.UI {
	id := "ha-svc-token"
	return app.Div().Class("col-md-9").
		Body(
			&FormControl{
				Compact: h.Compact,
				Content: func() []app.UI {
					return []app.UI{
						&FormText{
							ID:          id,
							Value:       h.happ.Name,
							BindTo:      &h.happ.Name,
							InputType:   "text",
							ReadOnly:    h.ReadOnly,
							Placeholder: "Not set",
						},
						&FormLabel{
							For:   id,
							Label: "Service Name",
						},
					}
				},
			},
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
				BindTo:   &h.happ.Active,
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

	r := &models.ApiResult{}
	err := frontend.PutHostedApp(u.String(), h.originalData.ID.Hex(), h.happ, &r)

	if err != nil {
		return
	}
	h.readonlyView(true)
	ctx.NewAction(actionHostedAppReload)
}
