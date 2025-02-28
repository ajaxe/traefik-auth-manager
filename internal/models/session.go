package models

import (
	"encoding/gob"

	"github.com/coreos/go-oidc/v3/oidc"
)

func init() {
	gob.Register(&Session{})
}

type Session struct {
	User SessionUser
}

type SessionUser struct {
	Name    string
	Sub     string
	Email   string
	Picture string
	IdpName string
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
