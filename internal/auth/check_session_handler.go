package auth

import (
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/db"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthCheckSession() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		sess, e := session.Get(sessionToken, c)
		if e != nil {
			err = echo.ErrUnauthorized
		}
		s := sess.Values[keyUserSession].(string)

		if s == "" {
			err = echo.ErrUnauthorized
		}

		existing, e := db.SessionByID(s)
		if e != nil || existing == nil {
			if e != nil {
				c.Logger().Errorf("failed to get session by id: %v: %v", s, e)
			} else {
				c.Logger().Errorf("invalid session id: %v", s)
			}
			err = echo.ErrUnauthorized
		}

		if err == echo.ErrUnauthorized {
			sess.Options.MaxAge = -1
			if e := sess.Save(c.Request(), c.Response()); e != nil {
				c.Logger().Errorf("faied to delete session cookie: %v", e)
			}
			return
		}

		err = c.JSON(http.StatusOK, existing)
		return
	}
}
