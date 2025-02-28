package auth

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthCheckSession() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		sess, err := session.Get(sessionToken, c)
		if err != nil {
			return echo.ErrUnauthorized
		}
		s := sess.Values[userSessionKey]
		return c.JSON(http.StatusOK, s)
	}
}
