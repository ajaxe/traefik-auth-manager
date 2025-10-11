package handlers

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type appUserHandler struct {
	logger echo.Logger
	dbFn   func() *db.AppUserDataAccess
}

func AddAppUserHandlers(e *echo.Group, l echo.Logger) {
	h := &appUserHandler{
		logger: l,
		dbFn:   db.NewAppUserDataAccess(),
	}
	e.GET("/app-users", h.Users())
	e.POST("/app-users", h.CreateUser())
	e.PUT(fmt.Sprintf("/app-users/%s", idParam.Param()), h.UpdateUser(idParam))
	e.DELETE(fmt.Sprintf("/app-users/%s", idParam.Param()), h.DeleteUser(idParam))

	e.DELETE(fmt.Sprintf("/app-users/%s/hosted-app/%s", idParam.Param(), appIDParam.Param()),
		h.removeUserApp(idParam, appIDParam))
	e.PUT(fmt.Sprintf("/app-users/%s/hosted-app/%s", idParam.Param(), appIDParam.Param()),
		h.assignUserApp(idParam, appIDParam))
}

func (h *appUserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := &models.AppUserChange{}

		if err := c.Bind(u); err != nil {
			return helpers.NewAppError(http.StatusBadRequest, "Bad data.", nil)
		}

		if u.UserName == "" {
			return helpers.NewAppError(http.StatusBadRequest, "Username is required.", nil)
		}

		err := h.validatePassword(u.Password, u.ConfirmPassword)
		if err != nil {
			return err
		}

		hash, err := helpers.GenerateHashUsingBase64URL(u.Password)
		if err != nil {
			return helpers.ErrAppGeneric(fmt.Errorf("error generating hash: %v", err))
		}

		da := h.dbFn()
		e, _ := da.AppUserByUsername(u.UserName)
		if e.UserName != "" { //err != mongo.ErrNoDocuments {
			return helpers.NewAppError(http.StatusBadRequest, fmt.Sprintf("User '%s' already exists.", u.UserName), nil)
		}

		id, err := da.InsertAppUser(&models.AppUser{
			UserName: u.UserName,
			Password: hash,
			Active:   true,
		})
		if err != nil {
			return helpers.ErrAppGeneric(fmt.Errorf("error saving user: %v", err))
		}

		return c.JSON(http.StatusOK, models.NewApiIDResult(id))
	}
}

func (h *appUserHandler) UpdateUser(idParam apiParam) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Param(idParam.String())

		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return helpers.ErrAppBadID(err)
		}

		var d models.AppUserChange
		if err := c.Bind(&d); err != nil {
			return helpers.ErrInvalidData(err)
		}
		d.ID = id

		if err := h.updatePassword(&d); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &models.ApiResult{
			Success: true,
		})
	}
}

func (h *appUserHandler) DeleteUser(idParam apiParam) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Param(idParam.String())

		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid ID.", err)
		}

		da := h.dbFn()
		err = da.DeleteAppUserByID(id)

		if err != nil {
			return helpers.ErrAppGeneric(err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (h *appUserHandler) Users() echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := h.userApps()
		if err != nil {
			return err
		}
		r := &models.AppUserListResult{
			ApiResult: models.ApiResult{
				Success: true,
			},
			Data: u,
		}
		return c.JSON(http.StatusOK, r)
	}
}

func (h *appUserHandler) userApps() (u []*models.AppUser, err error) {
	da := h.dbFn()
	u, err = da.AppUsers()
	if err != nil {
		return
	}
	a, err := db.HostedApplications()
	if err != nil {
		return
	}
	appMap := make(map[string]*models.HostedApplication)
	for _, v := range a {
		appMap[v.ID.Hex()] = v
	}
	for _, o := range u {
		for _, apps := range o.Applications {
			e, ok := appMap[apps.HostAppId.Hex()]
			if ok {
				apps.Name = e.Name
			} else {
				_ = apps.HostAppId.Hex()
			}
		}
	}
	return
}

func (h *appUserHandler) updatePassword(d *models.AppUserChange) error {

	err := h.validatePassword(d.Password, d.ConfirmPassword)
	if err != nil {
		return err
	}

	hash, err := helpers.GenerateHashUsingBase64URL(d.Password)

	if err != nil {
		return helpers.ErrAppGeneric(fmt.Errorf("error generating hash: %v", err))
	}
	da := h.dbFn()
	err = da.UpdatePassword(&models.AppUser{
		ID:       d.ID,
		Password: hash,
	})

	if err != nil {
		return helpers.ErrAppGeneric(fmt.Errorf("error saving password: %v", err))
	}

	return nil
}

func (h *appUserHandler) removeUserApp(idParam, appIdParam apiParam) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param(idParam.String())
		appId := c.Param(appIdParam.String())

		if id == "" {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid user id", nil)
		}

		if appId == "" {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid app id", nil)
		}
		da := h.dbFn()
		existing, err := da.AppUserByID(id)
		if err != nil || existing == nil {
			return helpers.NewAppError(http.StatusBadRequest, "User does not exist", err)
		}

		for i := range existing.Applications {
			if existing.Applications[i].HostAppId.Hex() == appId {
				existing.Applications = slices.Delete(existing.Applications, i, i+1)
				break
			}
		}

		err = da.UpdateUserHostedApps(existing)
		if err != nil {
			return helpers.ErrAppGeneric(fmt.Errorf("error updating user app: %v", err))
		}

		return c.JSON(http.StatusOK, &models.ApiResult{
			Success: true,
		})
	}
}
func (h *appUserHandler) assignUserApp(idParam, appIdParam apiParam) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param(idParam.String())
		appId := c.Param(appIdParam.String())

		if id == "" {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid user id", nil)
		}

		if appId == "" {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid app id", nil)
		}

		da := h.dbFn()
		existing, err := da.AppUserByID(id)
		if err != nil || existing == nil {
			return helpers.NewAppError(http.StatusBadRequest, "User does not exist", err)
		}

		has := false
		for i := range existing.Applications {
			if existing.Applications[i].HostAppId.Hex() == appId {
				has = true
				break
			}
		}

		if !has {
			h, err := db.HostedApplicationByID(appId)
			if h == nil || err != nil {
				return helpers.ErrAppGeneric(fmt.Errorf("non-existing hosted app: %s: %v", appId, err))
			}

			existing.Applications = append(existing.Applications,
				&models.ApplicationIdentifier{HostAppId: h.ID})

			err = da.UpdateUserHostedApps(existing)
			if err != nil {
				return helpers.ErrAppGeneric(fmt.Errorf("error updating user app: %v", err))
			}
		}

		return c.JSON(http.StatusOK, &models.ApiResult{
			Success: true,
		})
	}
}

func (h *appUserHandler) validatePassword(pwd, cpwd string) error {
	if pwd == "" {
		return helpers.NewAppError(http.StatusBadRequest, "Password is required", nil)
	}
	if cpwd != "" && pwd != cpwd {
		return helpers.NewAppError(http.StatusBadRequest, "Passwords do not match.", nil)
	}

	return nil
}
