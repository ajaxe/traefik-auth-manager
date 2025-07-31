//go:build !wasm

package main

import (
	"github.com/ajaxe/traefik-auth-manager/internal/backend"
	"github.com/labstack/echo/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type empty struct{ app.Compo }

func Frontend() {
	app.RouteWithRegexp("/.*", app.NewZeroComponentFactory(&empty{}))
}

func Backend(ah *app.Handler) {
	e := backend.NewBackendApi()

	e.GET("/*", func(c echo.Context) error {
		ah.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	backend.Start(e)
}
