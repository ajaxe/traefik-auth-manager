package components

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Modal struct {
	app.Compo
	containerID string
	Title       string
	Show        bool
	Content     func() []app.UI

	titleID string
}

func (m *Modal) Render() app.UI {
	if m.titleID == "" {
		m.titleID = fmt.Sprintf("m-title-%v", time.Now().UnixMicro())
	}

	return app.If(m.Show, func() app.UI {
		return app.Div().
			ID(m.containerID).
			Body(
				app.Div().Class("modal-backdrop fade show").Style("display", "block"),
				app.Div().Class("modal fade show").
					Style("display", "block").
					ID("editUserModal").
					TabIndex(-1).
					Aria("labelledby", m.titleID).
					Body(
						app.Div().Class("modal-dialog modal-dialog-centered").Body(
							app.Div().Class("modal-content").Body(
								app.Div().Class("modal-header").Body(
									app.H5().Class("modal-title").ID(m.titleID).Text(m.Title),
									app.Button().Class("btn-close").
										Aria("label", "Close").
										DataSet("bs-dismiss", "modal").
										OnClick(m.close),
								),
								app.Div().Class("modal-body").Body(
									m.Content()...,
								),
							),
						),
					),
			)
	}).Else(func() app.UI { return app.Div() })
}
func (m *Modal) close(ctx app.Context, e app.Event) {
	ctx.NewAction(actionUserEditClose)
}
