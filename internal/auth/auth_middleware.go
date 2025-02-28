package auth

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Authenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			sess, err := session.Get(sessionToken, c)
			if err != nil {
				c.Error(err)
			}
			if isauth, ok := sess.Values["isauth"].(bool); !ok || !isauth {
				return echo.ErrUnauthorized
			}
			if err = next(c); err != nil {
				c.Error(err)
			}
			return
		}
	}
}
