package backend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ajaxe/traefik-auth-manager/internal/auth"
	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/handlers"
	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
)

func NewBackendApi() *echo.Echo {
	appConfig := helpers.MustLoadDefaultAppConfig()

	e := echo.New()
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(session.Middleware(
		sessions.NewCookieStore([]byte(appConfig.Session.SessionKey))))

	handlers.AddHealtcheck(e)

	cfg := auth.InitAuthConfig(appConfig)

	e.POST("/login", auth.AuthLogin(cfg))
	e.POST("/signout", auth.AuthSignOut(cfg))
	e.POST(appConfig.OAuth.CallbackPath, auth.AuthCallback(cfg))
	e.GET(appConfig.OAuth.SignOutCallbackPath, auth.AuthSignOutCallback(cfg))

	a := e.Group("/api")
	a.Use(auth.Authenticated())
	a.GET("/check", auth.AuthCheckSession())

	handlers.AddAppUserHandlers(a)

	e.GET("/route-list", func(c echo.Context) error {
		return c.JSON(http.StatusOK, e.Routes())
	})

	return e
}

// Start echo server with graceful hanlding of process termination.
func Start(e *echo.Echo) {
	cfg := helpers.MustLoadDefaultAppConfig()
	addr := fmt.Sprintf(":%v", cfg.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		var err error
		if cfg.UseTLS() {
			e.Logger.Info("starting server with tls")
			err = e.StartTLS(addr, cfg.Server.CertFile, cfg.Server.KeyFile)
		} else {
			e.Logger.Info("starting server without tls")
			err = e.Start(addr)
		}
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Terminate(ctx); err != nil {
		e.Logger.Errorf("failed to terminate db connection: %v", err)
	}
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("failed to shutdown server: %v", err)
	}
}
