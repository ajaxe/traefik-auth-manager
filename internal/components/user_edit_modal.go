package components

import (
	"net/url"

	"github.com/ajaxe/traefik-auth-manager/internal/frontend"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserEditModal struct {
	app.Compo
	user            *models.AppUser
	show            bool
	password        string
	confirmPassword string
	formResult      *models.ApiResult
	op              string // 'add' or 'edit'. Defaults to 'edit'
}

func (u *UserEditModal) OnMount(ctx app.Context) {
	ctx.Handle(actionUserAdd, u.showAddModal)
	ctx.Handle(actionUserEdit, u.showModal)
	ctx.Handle(actionUserEditClose, u.closeModal)
}

func (u *UserEditModal) Render() app.UI {
	helpers.AppLogf("UserEditModal render: show: %v", u.show)
	if u.op == "" {
		u.op = "edit"
	}
	return &Modal{
		Title:   "Edit User",
		Content: u.form(),
		Show:    u.show,
	}
}
func (u *UserEditModal) setData(d *models.AppUser) {
	u.show = true
	u.user = d
	u.password = ""
	u.confirmPassword = ""
	u.formResult = nil
}
func (u *UserEditModal) showModal(ctx app.Context, a app.Action) {
	d, ok := a.Value.(*models.AppUser)
	helpers.AppLog(ok, d)

	u.show = true
	u.setData(d)

	ctx.Update()
}
func (u *UserEditModal) showAddModal(ctx app.Context, a app.Action) {
	u.show = true
	u.setData(&models.AppUser{})
	u.op = "add"

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
			app.Form().
				OnSubmit(u.formSubmit).
				Body(
					u.username(),
					u.passwordControl(),
					u.confirmPasswordControl(),
					u.saveBtn(),
				),
		}
	}
}
func (u *UserEditModal) username() app.UI {
	id := "username-ro"
	ro := u.op == "edit"

	return &FormControl{
		Content: func() []app.UI {
			return []app.UI{
				&FormText{
					ID:       id,
					Value:    u.user.UserName,
					ReadOnly: ro,
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
					ID:        id,
					Value:     u.password,
					BindTo:    &u.password,
					InputType: "password",
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
					ID:        id,
					Value:     u.confirmPassword,
					BindTo:    &u.confirmPassword,
					InputType: "password",
				},
				&FormLabel{
					For:   id,
					Label: "Confirm Password",
				},
			}
		},
	}
}
func (u *UserEditModal) saveBtn() app.UI {
	return app.Div().Body(
		app.Button().Class("btn btn-primary").Text("Save"),
		u.formResultMessage(),
	)
}

func (u *UserEditModal) formResultMessage() app.UI {
	if u.formResult == nil {
		return app.Span()
	}
	ico := "bi bi-check-circle text-success"
	co := "text-success"
	m := "Saved"

	if !u.formResult.Success {
		ico = "bi bi-x-circle text-danger"
		co = "text-danger"
		m = u.formResult.ErrorMessage
	}

	return app.Span().Class("ms-2 fw-bold").Body(
		app.I().Class(ico),
		app.Span().Class("ms-1 "+co).Text(m),
	)
}
func (u *UserEditModal) formSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	b := app.Window().URL()
	b.Path = ""

	helpers.AppLogf("UserEditModal form submitted: user id: %v", u.user.ID.Hex())

	var r models.ApiResult
	var err error
	if u.op == "edit" {
		r, err = u.updateUser(b)
	} else {
		r, err = u.addUser(b)
	}

	if err != nil {
		helpers.AppLogf("UserEditModal form submit error: %v", err)
		r = models.ApiResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}
	}
	u.formResult = &r
}
func (u *UserEditModal) updateUser(baseURL *url.URL) (r models.ApiResult, err error) {
	err = frontend.PutUser(u.user.ID.Hex(), baseURL.String(), &models.AppUserChange{
		AppUser: models.AppUser{
			ID: u.user.ID,
		},
		Password:        u.password,
		ConfirmPassword: u.confirmPassword,
	}, &r)

	return
}
func (u *UserEditModal) addUser(baseURL *url.URL) (r models.ApiResult, err error) {
	err = frontend.PostUser(baseURL.String(), &models.AppUserChange{
		AppUser: models.AppUser{
			UserName: u.user.UserName,
		},
		Password:        u.password,
		ConfirmPassword: u.confirmPassword,
	}, &r)

	return
}
