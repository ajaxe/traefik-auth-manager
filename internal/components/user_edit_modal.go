package components

import (
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserEditModal struct {
	app.Compo
	user *models.AppUser
	show bool
}

func (u *UserEditModal) OnMount(ctx app.Context) {
	ctx.Handle(actionUserEdit, u.showModal)
	ctx.Handle(actionUserEditClose, u.closeModal)
}

func (u *UserEditModal) Render() app.UI {
	helpers.AppLogf("UserEditModal render: show: %v", u.show)
	return &Modal{
		Title:   "Edit User",
		Content: u.form(),
		Show:    u.show,
	}
}
func (u *UserEditModal) showModal(ctx app.Context, a app.Action) {
	d, ok := a.Value.(*models.AppUser)
	helpers.AppLog(ok, d)
	u.show = true
	u.user = d
	ctx.Update()
}
func (u *UserEditModal) closeModal(ctx app.Context, a app.Action) {
	u.show = false
	u.user = nil
	ctx.Update()
}
func (u *UserEditModal) form() func() []app.UI {
	return func() []app.UI {
		return []app.UI{
			u.usernameReadOnly(),
			u.passwordControl(),
			u.confirmPasswordControl(),
		}
	}
}
func (u *UserEditModal) usernameReadOnly() app.UI {
	id := "username-ro"
	return &FormControl{
		Content: func() []app.UI {
			return []app.UI{
				&FormText{
					ID:       id,
					Value:    u.user.UserName,
					ReadOnly: true,
				},
				&FormLabel{
					For:   id,
					Label: "Username",
				},
			}
		},
	}
}
func (u *UserEditModal) passwordControl() app.UI {
	id := "user-pwd"
	return &FormControl{
		Content: func() []app.UI {
			return []app.UI{
				&FormText{
					ID:    id,
					Value: "",
				},
				&FormLabel{
					For:   id,
					Label: "Password",
				},
			}
		},
	}
}

func (u *UserEditModal) confirmPasswordControl() app.UI {
	id := "user-cnf-pwd"
	return &FormControl{
		Content: func() []app.UI {
			return []app.UI{
				&FormText{
					ID:    id,
					Value: "",
				},
				&FormLabel{
					For:   id,
					Label: "Confirm Password",
				},
			}
		},
	}
}
