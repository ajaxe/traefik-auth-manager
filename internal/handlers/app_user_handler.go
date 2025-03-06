package handlers

import (
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/labstack/echo/v4"
)

type appUserHandler struct {
}

func AddAppUserHandlers(e *echo.Group) {
	h := &appUserHandler{}
	e.GET("/app-users", h.Users())
}

func (h *appUserHandler) Users() echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := h.userApps()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
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
