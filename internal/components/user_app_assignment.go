package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserAppAssignment struct {
	app.Compo
	ID   string
	user *frontend.AppUserView
}

func (u *UserAppAssignment) Render() app.UI {

	return app.Div().
		ID(u.ID).
		Body(
			u.listApps()...,
		)
}
func (u *UserAppAssignment) listApps() []app.UI {
	l := []app.UI{}

	for _, k := range u.user.Apps() {

		btn := &UserAppAssignmentBtn{
			ID:       k.HostAppId.Hex(),
			selected: k.Selected,
			text:     k.Name,
			userId:   k.UserID,
		}
		l = append(l, btn)

	}
	return l
}

type UserAppAssignmentBtn struct {
	app.Compo
	text     string
	ID       string
	selected bool
	userId   string
}

func (u *UserAppAssignmentBtn) Render() app.UI {
	c := "btn-outline-secondary"
	if u.selected {
		c = "btn-primary"
	}

	return app.Button().
		DataSet("app-id", u.ID).
		DataSet("selected", u.selected).
		Class("me-1 mb-1 btn btn-sm " + c).
		Text(u.text).
		OnClick(u.click)
}
func (u *UserAppAssignmentBtn) click(ctx app.Context, e app.Event) {
	e.PreventDefault()

	frontend.NewAppContext(ctx).
		ToggleUserApp(u.userId, u.ID, u.selected,
			func() {
				ctx.Dispatch(func(ctx app.Context) {
					u.selected = !u.selected
				})
			})
}
