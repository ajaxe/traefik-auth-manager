package handlers

import (
	"fmt"
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/labstack/echo/v4"
)

type healtcheckStatus struct {
	Status string `json:"status"`
}

func AddHealtcheck(e *echo.Echo) {

	e.GET("/healthcheck/live", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &healtcheckStatus{
			Status: "alive",
		})
	})

	e.GET("/healthcheck/ready", func(c echo.Context) error {
		_, err := db.NewClient()
		if err != nil {
			return fmt.Errorf("failed connection to db: %v", err)
		}
		if err = db.Ping(c.Request().Context()); err != nil {
			return fmt.Errorf("failed to ping db: %v", err)
		}

		return c.JSON(http.StatusOK, &healtcheckStatus{
			Status: "ready",
		})
	})
}
