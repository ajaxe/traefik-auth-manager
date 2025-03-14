package handlers

import (
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/labstack/echo/v4"
)

type hostedAppHandler struct {
}

func AddHostedAppHandlers(e *echo.Group) {
	h := &hostedAppHandler{}

	e.GET("/hosted-apps", h.HostedApps())
}

func (h *hostedAppHandler) HostedApps() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		d, err := db.HostedApplications()
		if err != nil {
			return
		}

		return c.JSON(http.StatusOK, &models.HostedAppListResult{
			ApiResult: models.ApiResult{
				Success: true,
			},
			Data: d,
		})
	}
}
