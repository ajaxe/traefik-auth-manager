package backend

import (
	"context"
	"fmt"
	"log"
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
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var traceProvider *sdktrace.TracerProvider

func NewBackendApi() *echo.Echo {
	appConfig := helpers.MustLoadDefaultAppConfig()

	e := echo.New()

	e.Logger.SetLevel(elog.INFO)
	e.HTTPErrorHandler = handlers.AppErrorHandler()

	if appConfig.Server.IsDev {
		e.Debug = true
	} else {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	}

	if appConfig.Tracing.Enabled {
		var err error
		traceProvider, err = initTracer(appConfig.Tracing.ServiceName)
		if err != nil {
			log.Fatal(err)
		}
		e.Logger.Info("setting up otelecho.Middleware")
		e.Use(otelecho.Middleware(appConfig.Tracing.ServiceName,
			otelecho.WithTracerProvider(traceProvider),
		))
	}

	// e.Use(tracerMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(session.Middleware(
		sessions.NewCookieStore([]byte(appConfig.Session.SessionKey))))

	handlers.AddHealtcheck(e)

	cfg := auth.InitAuthConfig(appConfig)

	e.POST("/login", auth.AuthLogin(cfg))
	e.POST("/signout", auth.AuthSignOut(cfg))
	e.GET("/signout", auth.AuthSignOut(cfg))
	e.POST(appConfig.OAuth.CallbackPath, auth.AuthCallback(cfg))
	e.GET(appConfig.OAuth.SignOutCallbackPath, auth.AuthSignOutCallback(cfg))

	a := e.Group("/api")
	a.Use(auth.Authenticated())
	a.GET("/check", auth.AuthCheckSession())

	handlers.AddAppUserHandlers(a, e.Logger)

	handlers.AddHostedAppHandlers(a)

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
	if traceProvider != nil {
		if err := traceProvider.Shutdown(ctx); err != nil {
			e.Logger.Errorf("failed to shutdown trace provider: %v", err)
		}
	}
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("failed to shutdown server: %v", err)
	}
}

func initTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		log.Println("CRITICAL DEBUG: OTEL_EXPORTER_OTLP_ENDPOINT environment variable is NOT SET or is EMPTY.")
	} else {
		log.Printf("SUCCESS DEBUG: Found OTEL_EXPORTER_OTLP_ENDPOINT: %s", otelEndpoint)
	}

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
