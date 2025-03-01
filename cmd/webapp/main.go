package main

import (
	"log"

	"github.com/ajaxe/traefik-auth-manager/internal/backend"
	"github.com/ajaxe/traefik-auth-manager/internal/handlers"
	"github.com/ajaxe/traefik-auth-manager/internal/pages"
	"github.com/labstack/echo/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/", func() app.Composer { return &pages.HomePage{} })
	app.Route("/home", func() app.Composer { return &pages.HomePage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersPage{} })
	app.Route("/apps", func() app.Composer { return &pages.AppsPage{} })

	app.RunWhenOnBrowser()

	e := backend.NewBackendApi()

	e.GET("/*", func(c echo.Context) error {
		handlers.GoAppHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	backend.Start(e)

}
