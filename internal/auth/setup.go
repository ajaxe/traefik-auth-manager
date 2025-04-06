package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type appOAuthConfig struct {
	provider  *oidc.Provider
	config    oauth2.Config
	context   context.Context
	appConfig *helpers.AppConfig
}

const (
	tokenNonce       = "nonce"
	tokenState       = "state"
	tokenVerifier    = "verifier"
	sessionAuthSeq   = "auth-seq"
	sessionLogoutSeq = "logout-seq"
	// gorrilla session token name
	sessionToken   = "session-token"
	keyUserSession = "user-session"
	keyIDToken     = "idtoken"
	keyIsAuth      = "isauth"
)

func InitAuthConfig(appConfig helpers.AppConfig) appOAuthConfig {
	oauthCtx := context.Background()
	oidcProvider, err := oidc.NewProvider(oauthCtx, appConfig.OAuth.Authority)

	if err != nil {
		log.Fatalf("failed to create oidc provider: %v", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     appConfig.OAuth.ClientID,
		ClientSecret: appConfig.OAuth.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  appConfig.OAuthLoginRedirectURL(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return appOAuthConfig{
		oidcProvider,
		oauthConfig,
		oauthCtx,
		&appConfig,
	}
}

func RedirectToHome(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/home")
}

func authorizeUser(s models.Session) error {
	c := helpers.MustLoadDefaultAppConfig()
	valid := false

	l := strings.ToLower(s.User.Email)

	for _, r := range c.Session.AuthorizedEmails {
		if strings.ToLower(r) == l {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("unauthorized user: %v", s.User.Email)
	}

	return nil
}
