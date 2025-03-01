package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const oidcEndSessionEndpointKey = "end_session_endpoint"

func AuthSignOut(cfg appOAuthConfig) echo.HandlerFunc {
	oidcSignoutURL, e := getSignoutURL(cfg.appConfig.OAuth.Authority)
	return func(c echo.Context) error {
		if e != nil {
			return e
		}
		p := &map[string]string{
			"url": oidcSignoutURL,
		}
		return c.JSON(http.StatusOK, p)
	}
}

func getSignoutURL(issuer string) (string, error) {
	wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(wellKnown)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var doc = make(map[string]interface{})
	err = json.Unmarshal(b, &doc)

	if err != nil {
		return "", err
	}
	v, ok := doc[oidcEndSessionEndpointKey].(string)

	if !ok {
		return "", fmt.Errorf("oidc well-know config did not contain key '%s'", oidcEndSessionEndpointKey)
	}

	return v, nil
}
