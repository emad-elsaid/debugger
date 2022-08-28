package ui

import "gioui.org/widget"

type (
	Clickable = widget.Clickable
)

type Clickables map[string]*Clickable

func NewClickables() Clickables {
	return map[string]*Clickable{}
}

func (c Clickables) Get(id string) *Clickable {
	if btn, ok := c[id]; ok {
		return btn
	}

	btn := new(Clickable)
	c[id] = btn

	return btn
}
