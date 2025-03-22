package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type CardListItem struct {
	app.Compo
	title       string
	actionItems func() []app.UI
	content     func() []app.UI
}

func (c *CardListItem) Render() app.UI {
	if c.actionItems == nil {
		c.actionItems = func() []app.UI { return []app.UI{} }
	}
	if c.content == nil {
		c.content = func() []app.UI { return []app.UI{} }
	}

	a := []app.UI{
		app.Div().Class("me-auto").
			Style("padding-top", "5px").
			Body(
				app.Span().Class("h5").Text(c.title),
				app.I().Class("bi bi-arrow-right ms-2"),
			),
	}
	a = append(a, c.actionItems()...)

	t := []app.UI{app.Div().Class("card-title d-flex").Body(a...)}
	t = append(t, c.content()...)

	return app.Div().Class("card").Style("margin-bottom", "10px").
		Body(
			app.Div().Class("card-body").
				Body(t...),
		)
}
