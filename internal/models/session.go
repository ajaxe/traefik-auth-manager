//go:build !wasm

package models

import (
	"encoding/gob"

	"github.com/coreos/go-oidc/v3/oidc"
)

func init() {
	gob.Register(&Session{})
}

func NewSession(u *oidc.UserInfo) (s Session, err error) {
	c := make(map[string]interface{})

	err = u.Claims(&c)

	s = Session{
		User: SessionUser{
			Name:    c["name"].(string),
			Sub:     u.Subject,
			Email:   u.Email,
			Picture: c["picture"].(string),
			IdpName: c["private:idp"].(string),
		},
	}
	return
}