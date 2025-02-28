package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func randomString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func AuthLogin(cfg appOAuthConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(sessionAuthSeq, c)
		if err != nil {
			return err
		}

		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   5 * 60,
			HttpOnly: true,
		}

		state, err := randomString(16)
		if err != nil {
			return err
		}

		nonce, err := randomString(16)
		if err != nil {
			return err
		}

		verifier := oauth2.GenerateVerifier()

		sess.Values[tokenState] = state
		sess.Values[tokenVerifier] = verifier
		sess.Values[tokenNonce] = nonce

		if err = sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, cfg.config.AuthCodeURL(state,
			oidc.Nonce(nonce),
			oauth2.SetAuthURLParam("response_mode", "form_post"),
			oauth2.S256ChallengeOption(verifier),
			oauth2.SetAuthURLParam("prompt", "select_account"),
		))
	}
}
