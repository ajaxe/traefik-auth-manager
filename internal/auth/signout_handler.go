package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const oidcEndSessionEndpointKey = "end_session_endpoint"

func AuthSignOut(cfg appOAuthConfig) echo.HandlerFunc {
	oidcSignoutURL, e := getSignOutURL(cfg.appConfig.OAuth.Authority)
	return func(c echo.Context) (err error) {
		if e != nil {
			return RedirectToHome(c)
		}

		u, err := doExternalOIDCSignout(c, oidcSignoutURL, cfg)

		if err != nil {
			return
		}

		sess, err := session.Get(sessionToken, c)
		if err != nil {
			return
		}

		id, ok := sess.Values[keyUserSession].(string)

		if ok {
			err = db.DeleteSessionByID(id)

			if err != nil {
				return
			}
		}

		return c.Redirect(http.StatusFound, u)
	}
}

func doExternalOIDCSignout(c echo.Context, signoutURL string, cfg appOAuthConfig) (u string, err error) {
	sess, err := session.Get(sessionLogoutSeq, c)

	if err != nil {
		return
	}

	state, err := randomString(16)
	if err != nil {
		return
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	sess.Values[tokenState] = state

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return
	}

	u = buildSignoutURL(signoutURL,
		map[string]string{
			"post_logout_redirect_uri": cfg.appConfig.OAuthSignOutRedirectURL(),
			"state":                    state,
			"client_id":                cfg.appConfig.OAuth.ClientID,
		})

	c.Logger().Infof("redirecting for signout: %v", u)

	return
}

func buildSignoutURL(u string, q map[string]string) string {
	var buf bytes.Buffer
	buf.WriteString(u)

	val := url.Values{}

	for k, v := range q {
		val.Add(k, v)
	}

	buf.WriteString("?")

	buf.WriteString(val.Encode())
	return buf.String()
}

func getSignOutURL(issuer string) (v string, err error) {
	wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(wellKnown)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	var doc = make(map[string]interface{})
	err = json.Unmarshal(b, &doc)

	if err != nil {
		return
	}

	v, ok := doc[oidcEndSessionEndpointKey].(string)

	if !ok {
		err = fmt.Errorf("oidc well-know config did not contain key '%s'", oidcEndSessionEndpointKey)
	}

	return
}
