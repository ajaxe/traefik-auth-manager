package components

import (
	"fmt"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserDeleteBtn struct {
	app.Compo
	user *frontend.AppUserView
	id   string
}

func (u *UserDeleteBtn) Render() app.UI {
	u.id = fmt.Sprintf("u-del-btn-%s", u.user.ID.Hex())
	return &DeleteBtn{
		ID: u.id,
		onDelete: func(ctx app.Context) {
			frontend.NewAppContext(ctx).RemoveUser(u.user.ID.Hex())
		},
	}
}
