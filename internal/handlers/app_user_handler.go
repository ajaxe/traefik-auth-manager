package handlers

import (
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type appUserHandler struct {
	logger echo.Logger
}

func AddAppUserHandlers(e *echo.Group, l echo.Logger) {
	h := &appUserHandler{
		logger: l,
	}
	e.GET("/app-users", h.Users())
	e.POST("/app-users", h.CreateUser())
	e.PUT("/app-users/:id", h.UpdateUser())
	e.DELETE("/app-users/:id", h.DeleteUser())
}

func (h *appUserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	}
}

func (h *appUserHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		i := c.Param("id")

		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid ID.", err)
		}

		var d models.AppUserChange
		if err := c.Bind(&d); err != nil {
			return helpers.NewAppError(http.StatusBadRequest, "Invalid data.", err)
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

func (h *appUserHandler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
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
	u, err = db.AppUsers()
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
			apps.Name = appMap[apps.HostAppId.Hex()].Name
		}
	}
	return
}
func (h *appUserHandler) updatePassword(d *models.AppUserChange) error {

	if d.Password == "" {
		return helpers.NewAppError(http.StatusBadRequest, "Password is required", nil)
	}
	if d.Password != d.ConfirmPassword {
		return helpers.NewAppError(http.StatusBadRequest, "Passwords do not match.", nil)
	}
	hash, err := helpers.GenerateHashUsingBase64URL(d.Password)

	if err != nil {
		h.logger.Errorf("error generating hash: %v", err)
		return helpers.ErrAppGeneric
	}

	err = db.UpdatePassword(&models.AppUser{
		ID:       d.ID,
		Password: hash,
	})

	if err != nil {
		h.logger.Errorf("error saving password: %v", err)
		return helpers.ErrAppGeneric
	}

	return nil
}
