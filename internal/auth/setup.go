package auth

import (
	"context"
	"log"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type appOAuthConfig struct {
	provider  *oidc.Provider
	config    oauth2.Config
	context   context.Context
	appConfig *helpers.AppConfig
}

const (
	tokenNonce     = "nonce"
	tokenState     = "state"
	tokenVerifier  = "verifier"
	sessionAuthSeq = "auth-seq"
	// gorrilla session token name
	sessionToken   = "session-token"
	userSessionKey = "user-session"
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
		RedirectURL:  appConfig.OAuthRedirectURL(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return appOAuthConfig{
		oidcProvider,
		oauthConfig,
		oauthCtx,
		&appConfig,
	}
}
