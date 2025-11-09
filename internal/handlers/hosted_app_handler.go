package handlers

import (
	"context"
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
	logger echo.Logger
	dbFn   func(context.Context) *db.HostedAppDataAccess
}

func AddHostedAppHandlers(e *echo.Group, l echo.Logger) {
	h := &hostedAppHandler{
		logger: l,
		dbFn:   db.NewHostedAppDataAccess(),
	}

	e.GET("/hosted-apps", h.HostedApps())
	e.POST("/hosted-apps", h.CreateHostedApps())
	e.PUT(fmt.Sprintf("/hosted-apps/%s", idParam.Param()), h.UpdateHostedApp(idParam))
	e.DELETE(fmt.Sprintf("/hosted-apps/%s", idParam.Param()), h.DeleteHostedApp(idParam))
}

func (h *hostedAppHandler) HostedApps() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		d, err := h.dbFn(c.Request().Context()).HostedApplications()
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
		var d models.HostedApplication
		if err := c.Bind(&d); err != nil {
			return helpers.ErrInvalidData(err)
		}

		err = h.validate(d)
		if err != nil {
			return err
		}

		db := h.dbFn(c.Request().Context())

		existing, _ := db.HostedApplicationByServiceToken(d.ServiceToken)
		if existing != nil && existing.ServiceToken == d.ServiceToken {
			return helpers.NewAppError(http.StatusBadRequest, "Service token already exists.", nil)
		}

		id, err := db.InsertHostedApplication(&d)
		if err != nil {
			return helpers.ErrAppGeneric(err)
		}

		return c.JSON(http.StatusOK, models.NewApiIDResult(id))
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

		err = h.dbFn(c.Request().Context()).UpdateHostedApplication(&d)
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
		return helpers.ErrAppRequired("Service name")
	}
	if m.ServiceToken == "" {
		return helpers.ErrAppRequired("Service token")
	}
	if m.ServiceURL == "" {
		return helpers.ErrAppRequired("Service URL")
	}
	_, err := url.ParseRequestURI(m.ServiceURL)
	if err != nil {
		return helpers.NewAppError(http.StatusBadRequest, "Invalid service URL.", err)
	}

	return nil
}
func (h *hostedAppHandler) DeleteHostedApp(p apiParam) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		i := c.Param(idParam.String())

		id, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return helpers.ErrAppBadID(err)
		}

		if id.IsZero() {
			return helpers.ErrAppRequired("Host app ID")
		}

		app, err := h.dbFn(c.Request().Context()).HostedApplicationByID(i)
		if err != nil {
			return helpers.ErrAppGeneric(err)
		}

		err = db.DeleteHostedAppByID(id)
		if err != nil {
			return helpers.ErrAppGeneric(err)
		}

		return c.JSON(http.StatusAccepted, &models.HostedAppListResult{
			ApiResult: models.ApiResult{
				Success: true,
			},
			Data: []*models.HostedApplication{app},
		})
	}
}
