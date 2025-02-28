package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthCheckSession() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(http.StatusOK, "ok")
	}
}
