package handlers

import (
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/labstack/echo/v4"
)

func AppErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		c.Logger().Errorf("error: %v", err)
		st := http.StatusInternalServerError
		m := err.Error()
		if ae, ok := err.(*helpers.AppError); ok {
			st = ae.HTTPStatus()
			m = ae.Message()
		}
		c.JSON(st, &models.ApiResult{
			Success:      false,
			ErrorMessage: m,
		})
	}
}
