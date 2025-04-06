package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type hostedAppHandler struct {
}

func AddHostedAppHandlers(e *echo.Group) {
	h := &hostedAppHandler{}

	e.GET("/hosted-apps", h.HostedApps())
	e.POST("/hosted-apps", h.CreateHostedApps())
	e.PUT(fmt.Sprintf("/hosted-apps/%s", idParam.Param()), h.UpdateHostedApp(idParam))
}

func (h *hostedAppHandler) HostedApps() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		d, err := db.HostedApplications()
		if err != nil {
			return
		}

		sort.Slice(d, func(i, j int) bool {
			return strings.ToLower(d[i].Name) < strings.ToLower(d[j].Name)
		})

		return c.JSON(http.StatusOK, &models.HostedAppListResult{
			ApiResult: models.ApiResult{
				Success: true,
			},
			Data: d,
		})
	}
}
func (h *hostedAppHandler) CreateHostedApps() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(http.StatusOK, "Noop")
	}
}
func (h *hostedAppHandler) UpdateHostedApp(p apiParam) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		i := c.Param(idParam.String())

		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return helpers.ErrAppBadID(err)
		}

		var d models.HostedApplication
		if err := c.Bind(&d); err != nil {
			return helpers.ErrInvalidData(err)
		}
		d.ID = id

		if d.ID.IsZero() {
			return helpers.ErrAppRequired("Host app ID")
		}

		err = h.validate(d)
		if err != nil {
			return err
		}

		err = db.UpdateHostedApplication(&d)
		if err != nil {
			return helpers.ErrAppGeneric(err)
		}

		return c.JSON(http.StatusAccepted, &models.ApiResult{
			Success: true,
		})
	}
}
func (h *hostedAppHandler) validate(m models.HostedApplication) error {
	if m.Name == "" {
		return helpers.ErrAppRequired("Hosted app name")
	}
	if m.ServiceToken == "" {
		return helpers.ErrAppRequired("Hosted app service token")
	}
	if m.ServiceURL == "" {
		return helpers.ErrAppRequired("Hosted app service URL")
	}
	if _, err := url.Parse(m.ServiceURL); err != nil {
		return helpers.NewAppError(http.StatusBadRequest, "Invalid service URL.", err)
	}
	return nil
}
