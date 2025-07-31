// single entry point for the webapp & backend api
package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	Frontend()

	app.RunWhenOnBrowser()

	Backend(GoAppHandler)
}
