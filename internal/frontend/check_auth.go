package frontend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ajaxe/traefik-auth-manager/internal/models"
)

func CheckAuth(u string) (models.Session, error) {
	var s models.Session
	err := httpGet(fmt.Sprintf("%s/api/check", u), &s)

	if err != nil {
		return models.Session{}, err
	}

	return s, nil
}

func httpGet(u string, v interface{}) error {
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error code: %v", res.StatusCode)
	}

	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &v)

	return err
}
