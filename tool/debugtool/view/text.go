package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type TextView struct {
	entry   *widget.Entry
	msgChan chan string
	raw     string
	btn     *widget.Button
}

func (t *TextView) Chan() chan string {
	return t.msgChan
}
func (t *TextView) Daemon() {
	go func() {
		for msg := range t.msgChan {
			t.raw += msg + "\n"
			t.entry.SetText(t.raw)
		}
	}()
}
func (t *TextView) Join() *fyne.Container {
	//text := container.NewVBox(t.entry)
	//text := container.New(layout.NewVBoxLayout(), t.entry)
	content := container.New(
		layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewBorder(t.btn, nil, nil, nil, t.entry),
	)
	return content
}
func (t *TextView) Clear() {
	t.raw = ""
	t.entry.SetText("")
}
func NewTextView() *TextView {
	t := &TextView{
		entry:   widget.NewMultiLineEntry(),
		msgChan: make(chan string, 1024),
	}
	t.btn = widget.NewButton("clear", t.Clear)
	return t
}
