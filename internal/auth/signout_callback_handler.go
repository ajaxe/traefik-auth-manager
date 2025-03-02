package auth

import (
	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthSignOutCallback(cfg appOAuthConfig) echo.HandlerFunc {
	return func(c echo.Context) (e error) {
		defer func() { e = RedirectToHome(c) }()

		sess, err := session.Get(sessionLogoutSeq, c)

		if err != nil {
			c.Logger().Errorf("error getting logout session: v%", err)
			return
		}

		s := c.Request().URL.Query()["state"]
		state := sess.Values[tokenState].(string)

		if len(s) != 1 && s[0] != state {
			return
		}

		usess, err := session.Get(sessionToken, c)

		if err != nil {
			c.Logger().Errorf("error getting logout session: v%", err)
			return
		}
		defer func() {
			if err = usess.Save(c.Request(), c.Response()); err != nil {
				c.Logger().Errorf("error removing user session: v%", err)
			}
		}()

		usess.Options = &sessions.Options{
			MaxAge: -1,
		}
		id := usess.Values[keyUserSession].(string)

		err = db.DeleteSessionByID(id)

		return
	}
}
