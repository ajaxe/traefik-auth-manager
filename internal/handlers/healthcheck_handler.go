package handlers

import (
	"fmt"
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/labstack/echo/v4"
)

func AddHealtcheck(e *echo.Echo) {
	e.GET("/healthcheck", func(c echo.Context) error {
		_, err := db.NewClient()
		if err != nil {
			return fmt.Errorf("failed connection to db: %v", err)
		}
		if err = db.Ping(); err != nil {
			return fmt.Errorf("failed to ping db: %v", err)
		}

		return c.String(http.StatusOK, "OK")
	})
}
