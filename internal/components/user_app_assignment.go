package components

import (
	"sort"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserAppAssignment struct {
	app.Compo
	ID       string
	userApps []*models.ApplicationIdentifier
	allApps  map[string]*models.HostedApplication
	userId   string
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
	a := make(map[string]bool)

	for _, r := range u.userApps {
		a[strings.ToLower(r.Name)] = true
	}

	keys := []string{}
	for k := range u.allApps {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		_, ok := a[k]

		btn := &UserAppAssignmentBtn{
			ID:       u.allApps[k].ID.Hex(),
			selected: ok,
			text:     u.allApps[k].Name,
			userId:   u.userId,
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
		Class("me-1 btn btn-sm " + c).
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
