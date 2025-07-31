//go:build wasm

package main

import (
	"github.com/ajaxe/traefik-auth-manager/internal/pages"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func Frontend() {
	app.Route("/", func() app.Composer { return &pages.HomePage{} })
	app.Route("/home", func() app.Composer { return &pages.HomePage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersPage{} })
	app.Route("/apps", func() app.Composer { return &pages.AppsPage{} })
}

func Backend(_ *app.Handler) {

}
