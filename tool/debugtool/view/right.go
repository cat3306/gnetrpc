package view

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cat3306/gnetrpc/protocol"
)

var (
	CodeSelectList = []string{"json", "string"}
)

type SendView struct {
	Title         *widget.Label
	Method        *Input
	Path          *Input
	SendBtn       *widget.Button
	BodyInput     *Input
	CodeSelect    *Select
	Metadata      *Input
	serializeType protocol.SerializeType
}

func (s *SendView) Join() *fyne.Container {
	box := container.NewVBox(
		s.Title,
		s.Path.Join(),
		s.Method.Join(),
		s.CodeSelect.Join(),
		s.Metadata.Join(),
		s.BodyInput.Join(),
		s.SendBtn,
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
			metaData := make(map[string]string)
			if s.Metadata.Entry.Text != "" {
				err := json.Unmarshal([]byte(s.Metadata.Entry.Text), &metaData)
				if err != nil {
					GlobalText.msgChan <- err.Error()
				}
			}

			RpcClient.Client.Call(s.Path.Entry.Text, s.Method.Entry.Text, metaData, s.serializeType, body)
		}
	})

	bodyI := &Input{
		Layout:    layout.NewFormLayout(),
		Label:     widget.NewLabel("body:"),
		Entry:     widget.NewMultiLineEntry(),
		Container: nil,
	}
	s.BodyInput = bodyI

	tmpSelect := widget.NewSelect(CodeSelectList, func(tmp string) {
		if tmp == "json" {
			s.serializeType = protocol.Json
		} else if tmp == "string" {
			s.serializeType = protocol.String
		} else {
			s.serializeType = protocol.Json
		}
	})
	tmpSelect.SetSelectedIndex(0)

	s.CodeSelect = &Select{
		Layout: layout.NewFormLayout(),
		Label:  widget.NewLabel("serialize type:"),
		Select: tmpSelect,
	}

	metadata := &Input{
		Layout:    layout.NewFormLayout(),
		Label:     widget.NewLabel("metadata:"),
		Entry:     widget.NewMultiLineEntry(),
		Container: nil,
	}
	s.Metadata = metadata
	return s.Join()
}