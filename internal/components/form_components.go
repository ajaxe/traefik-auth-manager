package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type FormLabel struct {
	app.Compo
	For   string
	Label string
}

func (l *FormLabel) Render() app.UI {
	return app.Label().For(l.For).
		Text(l.Label)
}

type FormText struct {
	app.Compo
	ID       string
	Value    string
	ReadOnly bool
}

func (t *FormText) Render() app.UI {
	c := "form-control"
	if t.ReadOnly {
		c += "-plaintext"
	}
	return app.Input().
		Type("text").
		ReadOnly(t.ReadOnly).
		Placeholder(t.Value).
		Class(c).
		ID(t.ID).
		Value(t.Value)
}

type FormControl struct {
	app.Compo
	Content func() []app.UI
}

func (f *FormControl) Render() app.UI {
	return app.Div().Class("form-floating mb-3").Body(f.Content()...)
}
