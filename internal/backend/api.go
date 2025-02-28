package backend

import (
	"github.com/ajaxe/traefik-auth-manager/internal/auth"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewBackendApi() *echo.Echo {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	appConfig := helpers.MustLoadDefaultAppConfig()

	e.Use(session.Middleware(
		sessions.NewCookieStore([]byte(appConfig.Session.SessionKey))))

	cfg := auth.InitAuthConfig(appConfig)
	e.GET("/login", auth.AuthLogin(cfg)) // for testing only
	e.POST("/login", auth.AuthLogin(cfg))
	e.POST(appConfig.OAuth.CallbackPath, auth.AuthCallback(cfg))

	a := e.Group("/")
	a.Use(auth.Authenticated())
	a.GET("/check", auth.AuthCheckSession())

	return e
}
