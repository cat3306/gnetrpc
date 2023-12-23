package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cat3306/gnetrpc/protocol"
)

type SendView struct {
	Title         *widget.Label
	Method        *Input
	Path          *Input
	SendBtn       *widget.Button
	BodyInput     *Input
	CodeSelect    *widget.Select
	serializeType protocol.SerializeType
}

func (s *SendView) Join() *fyne.Container {
	box := container.NewVBox(
		s.Title,
		s.Path.Join(),
		s.Method.Join(),
		s.CodeSelect,
		s.SendBtn,
		s.BodyInput.Join(),
	)

	border := container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), box)
	return border
}

func (s *SendView) Init() *fyne.Container {
	ie := widget.NewEntry()
	ie.SetText("Builtin")
	s.Path = &Input{
		Layout: layout.NewFormLayout(),
		Label:  widget.NewLabel("service path:"),
		Entry:  ie,
	}

	me := widget.NewEntry()
	me.SetText("Heartbeat")
	s.Method = &Input{
		Layout: layout.NewFormLayout(),
		Label:  widget.NewLabel("method:"),
		Entry:  me,
	}
	s.Title = widget.NewLabel("send binary")
	s.SendBtn = widget.NewButton("send", func() {
		if RpcClient.Client != nil {
			var body interface{}
			if s.serializeType == protocol.Json {
				body = []byte(s.BodyInput.Entry.Text)
			} else if s.serializeType == protocol.String {
				body = s.BodyInput.Entry.Text
			}
			RpcClient.Client.Call(s.Path.Entry.Text, s.Method.Entry.Text, nil, s.serializeType, body)
		}
	})

	bodyI := &Input{
		Layout:    layout.NewFormLayout(),
		Label:     widget.NewLabel("body:"),
		Entry:     widget.NewMultiLineEntry(),
		Container: nil,
	}
	s.BodyInput = bodyI

	s.CodeSelect = widget.NewSelect([]string{"json", "string"}, func(tmp string) {
		if tmp == "json" {
			s.serializeType = protocol.Json
		} else if tmp == "string" {
			s.serializeType = protocol.String
		} else {
			s.serializeType = protocol.Json
		}
	})
	return s.Join()
}
