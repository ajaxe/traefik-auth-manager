package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

func AppCodeUpdate() *CodeUpdate {
	return &CodeUpdate{}
}
func AppNavBar() *NavBar {
	return &NavBar{}
}

func AppNavBarItems(options NavListOptions) *NavBarItems {
	return &NavBarItems{
		itemTextColor: options.TextColor,
		listCSS:       options.ListCSS,
		listParent:    app.Ul().Class("navbar-nav"),
	}
}
func AppSignoutBtn() *SignOutBtn {
	return &SignOutBtn{}
}