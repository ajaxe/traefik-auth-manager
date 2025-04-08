package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type HostedAppDeleteBtn struct {
	app.Compo
	Happ *models.HostedApplication
}

func (h *HostedAppDeleteBtn) Render() app.UI {
	return &DeleteBtn{
		ID: h.Happ.ID.Hex(),
		onDelete: func(ctx app.Context) {
			frontend.NewAppContext(ctx).RemoveHostedApp(h.Happ.ID.Hex())
		},
	}
}
