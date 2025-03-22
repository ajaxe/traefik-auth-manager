package components

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

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
	ID        string
	Value     string
	ReadOnly  bool
	BindTo    any
	InputType string
}

func (t *FormText) Render() app.UI {
	c := "form-control"
	if t.ReadOnly {
		c += "-plaintext"
	}
	it := t.InputType
	if it == "" {
		it = "text"
	}

	elem := app.Input().
		Type(it).
		ReadOnly(t.ReadOnly).
		Placeholder(t.Value).
		Class(c).
		ID(t.ID).
		Value(t.Value)

	if !t.ReadOnly {
		elem.OnChange(t.ValueTo(t.BindTo))
	}

	return elem
}

type FormControl struct {
	app.Compo
	Content func() []app.UI
}

func (f *FormControl) Render() app.UI {
	return app.Div().Class("form-floating mb-3").Body(f.Content()...)
}

type FormCheckbox struct {
	app.Compo
	role   string // empty or "switch"
	bindTo any
	value  bool
	label  string
}

func (c *FormCheckbox) Render() app.UI {
	s := ""
	if c.role == "switch" {
		s = "form-switch"
	}

	id := fmt.Sprintf("chk-%v", time.Now().UnixMicro())

	input := app.Input().
		Class("form-check-input").
		Type("checkbox").
		Value(c.value).
		Checked(c.value).
		ID(id)

	if c.bindTo != nil {
		input.OnChange(c.ValueTo(c.bindTo))
	}

	return app.Div().Class("form-check "+s).
		Body(
			input,
			app.Label().Class("form-check-label").
				For(id).
				Text(c.label),
		)
}
