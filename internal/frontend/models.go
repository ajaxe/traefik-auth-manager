package frontend

import (
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
)

type HostedAppCollection struct {
	apps           []*models.HostedApplication
	appMapInternal map[string]*models.HostedApplication
}

func (h *HostedAppCollection) Names() []string {
	l := []string{}
	for _, v := range h.apps {
		l = append(l, v.Name)
	}
	return l
}
func (h *HostedAppCollection) Lookup(name string) *models.HostedApplication {
	if h.appMapInternal == nil {
		m := make(map[string]*models.HostedApplication)

		for _, k := range h.apps {
			m[strings.ToLower(k.Name)] = k
		}

		h.appMapInternal = m
	}
	return h.appMapInternal[strings.ToLower(name)]
}

type AppUserView struct {
	models.AppUser
	allApps *HostedAppCollection
}

func (a *AppUserView) Apps() []UserApplicationView {
	r := []UserApplicationView{}
	m := make(map[string]bool)

	for _, r := range a.Applications {
		m[r.Name] = true
	}

	for _, k := range a.allApps.Names() {
		_, ok := m[k]

		btn := UserApplicationView{
			ApplicationIdentifier: models.ApplicationIdentifier{
				HostAppId: a.allApps.Lookup(k).ID,
				Name:      a.allApps.Lookup(k).Name,
			},
			Selected: ok,
			UserID:   a.ID.Hex(),
		}
		r = append(r, btn)
	}

	return r
}

type UserApplicationView struct {
	models.ApplicationIdentifier
	Selected bool
	UserID   string
}

func NewUserListViewData(user []*models.AppUser, apps []*models.HostedApplication) UserListViewData {
	s := &HostedAppCollection{apps: apps}
	l := UserListViewData{
		Users: make([]*AppUserView, len(user)),
	}
	for i, r := range user {
		l.Users[i] = &AppUserView{
			AppUser: *r,
			allApps: s,
		}
	}
	return l
}

type UserListViewData struct {
	Users []*AppUserView
}
