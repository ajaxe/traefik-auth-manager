package main

import (
	"log"

	"github.com/ajaxe/traefik-auth-manager/internal/handlers"
	"github.com/ajaxe/traefik-auth-manager/internal/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/", func() app.Composer { return &pages.HomePage{} })
	app.Route("/home", func() app.Composer { return &pages.HomePage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersPage{} })
	app.Route("/apps", func() app.Composer { return &pages.AppsPage{} })

	app.RunWhenOnBrowser()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/*", func(c echo.Context) error {
		handlers.GoAppHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
