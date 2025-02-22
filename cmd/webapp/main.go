package main

import (
	"log"

	"github.com/ajaxe/traefik-auth-manager/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type hello struct {
	app.Compo
	name string
}

func (h *hello) Render() app.UI {
	return app.Map().Body(
		app.H1().Text("Hello, "+h.name),
		app.Input().Type("text").Value(h.name).OnChange(h.OnChange),
	)
}

func (h *hello) OnChange(ctx app.Context, e app.Event) {
	h.name = ctx.JSSrc().Get("value").String()
	ctx.Update()
}

func main() {
	app.Route("/", func() app.Composer { return &hello{name: "World!"} })

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
