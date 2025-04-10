package components

import (
	"fmt"
	"time"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserListItem struct {
	app.Compo
	user *frontend.AppUserView
}

func (ul *UserListItem) Render() app.UI {
	i := fmt.Sprintf("c%v", time.Now().UnixMilli())
	return &CardListItem{
		Title: ul.user.UserName,
		actionItems: []app.UI{
			&UserEditBtn{user: ul.user},
			&UserDeleteBtn{user: ul.user},
		},
		content: []app.UI{
			&UserAppAssignment{
				ID:   i,
				user: ul.user,
			},
		},
	}
}
