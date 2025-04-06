package auth

import (
	"fmt"
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// Handles authentication callback when returning from the identity provider.
// Validates state, nonce & exchanges code for access_token, id_token
func AuthCallback(cfg appOAuthConfig) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		sess, err := session.Get(sessionAuthSeq, c)
		if err != nil {
			c.Logger().Errorf("failed to get session: %v", err)
			return echo.ErrInternalServerError
		}
		sess.Options.MaxAge = -1
		err = sess.Save(c.Request(), c.Response())

		if err != nil {
			c.Logger().Errorf("failed to remove %s session: %v", sessionAuthSeq, err)
			return echo.ErrInternalServerError
		}

		expectedState := sess.Values[tokenState]
		incomingState := c.FormValue(tokenState)

		if expectedState != incomingState {
			c.Logger().
				Errorf("oauth state does not match, expected - %s, received - %s",
					expectedState, incomingState)
			return echo.ErrBadRequest
		}

		verifier := fmt.Sprintf("%s", sess.Values[tokenVerifier])
		token, err := cfg.config.Exchange(cfg.context, c.FormValue("code"),
			oauth2.VerifierOption(verifier))

		if err != nil {
			c.Logger().
				Errorf("failed to exchange token for access_token: %v", err)
			return echo.ErrBadRequest
		}

		idtoken, err := validatedIDToken(token, c, sess, cfg)
		if err != nil {
			return err
		}

		userInfo, err := cfg.provider.UserInfo(cfg.context, oauth2.StaticTokenSource(token))
		if err != nil {
			c.Logger().Errorf("failed to get userinfo: %v", err)
			return echo.ErrInternalServerError
		}

		s, err := models.NewSession(userInfo)

		if err != nil {
			c.Logger().Errorf("failed to create user session: %v", err)
			return echo.ErrInternalServerError
		}

		tokenSess, err := session.Get(sessionToken, c)
		if err != nil {
			return err
		}

		//t := int(token.Expiry.Sub(time.Now().UTC()).Seconds())
		tokenSess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
			Secure:   true,
		}

		err = authorizeUser(s)
		if err != nil {
			c.Logger().Errorf("failed to authorize user: %v", err)
			return c.Redirect(http.StatusFound, "/signout") //echo.ErrForbidden
		}

		id, err := db.InsertSession(s)
		if err != nil {
			c.Logger().Errorf("failed to create db user session: %v", err)
			return echo.ErrInternalServerError
		}

		tokenSess.Values[keyIsAuth] = true
		tokenSess.Values[keyUserSession] = id
		tokenSess.Values[keyIDToken] = idtoken

		err = tokenSess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, "/home")
	}
}

func validatedIDToken(token *oauth2.Token, c echo.Context,
	store *sessions.Session,
	cfg appOAuthConfig) (rawIDToken string, err error) {

	verifier := cfg.provider.Verifier(&oidc.Config{
		ClientID: cfg.config.ClientID,
	})

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.Logger().Error("no id_token field in oauth2 token.")
		err = echo.ErrInternalServerError
		return
	}
	idToken, err := verifier.Verify(cfg.context, rawIDToken)
	if err != nil {
		c.Logger().Errorf("failed to verify id_token: %v", err)
		err = echo.ErrInternalServerError
		return
	}

	nonce, ok := store.Values["nonce"].(string)
	if !ok {
		c.Logger().Error("nonce not found")
		err = echo.ErrBadRequest
	}
	if idToken.Nonce != nonce {
		c.Logger().Error("id_token nonce did not match")
		err = echo.ErrBadRequest
	}

	return
}
