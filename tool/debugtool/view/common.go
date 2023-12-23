package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Input struct {
	Layout    fyne.Layout
	Label     *widget.Label
	Entry     *widget.Entry
	Container *fyne.Container
}

func (i *Input) Join() *fyne.Container {
	if i.Container == nil {
		i.Container = container.New(i.Layout, i.Label, i.Entry)
	}
	return i.Container
}
