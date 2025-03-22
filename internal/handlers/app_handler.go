package handlers

import "github.com/maxence-charriere/go-app/v10/pkg/app"

var GoAppHandler = &app.Handler{
	Name:        "Proxy Auth Manager",
	Title:       "Proxy Auth Manager",
	Description: "Helper application to manage forward authentication for Traefik.",
	Icon:        app.Icon{Default: "/web/favicon.ico", SVG: "/web/favicon.ico"},

	Styles: []string{"/web/css/bootstrap.min.css", "/web/css/common.css", "/web/font/bootstrap-icons.min.css"},
	Scripts: []string{"/web/scripts/bootstrap.bundle.min.js",
		"/web/scripts/popper.min.js",
		"/web/scripts/cash.min.js",
		"/web/scripts/common.js",
	},
	Fonts: []string{"/web/font/fonts/bootstrap-icons.woff2"},
}
