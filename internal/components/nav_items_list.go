package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type NavBarItems struct {
	app.Compo
	listParent    app.HTMLUl
	itemTextColor string
	listCSS       string
}

func (h *NavBarItems) Render() app.UI {
	if h.listCSS == "" {
		h.listCSS = "navbar-nav mr-auto"
	}
	return h.listParent.Class(h.listCSS).Body(h.items()...)
}
func (h *NavBarItems) items() []app.UI {
	return []app.UI{
		h.item("/home", "Home"),
		h.item("/users", "Users"),
		h.item("/apps", "Applications"),
	}
}
func (h *NavBarItems) item(href, text string) app.UI {
	return app.Li().Class("nav-item").Body(
		app.A().Class("nav-link " + h.itemTextColor).Href(href).Text(text),
	)
}
