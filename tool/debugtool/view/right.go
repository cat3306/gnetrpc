package view

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cat3306/gnetrpc/protocol"
)

var (
	CodeSelectList = []string{"json", "string", "none"}
)

type SendView struct {
	Title         *widget.Label
	Method        *Input
	Path          *Input
	SendBtn       *widget.Button
	BodyInput     *Input
	CodeSelect    *Select
	Metadata      *Input
	Checkbox      *widget.Check
	SendInterval  *Input
	serializeType protocol.SerializeType
}

func (s *SendView) Join() *fyne.Container {
	btn := container.NewHBox(s.Checkbox, s.SendInterval.Join(), s.SendBtn)
	box := container.NewVBox(

		s.Path.Join(),
		s.Method.Join(),
		s.CodeSelect.Join(),
		s.Metadata.Join(),
		s.BodyInput.Join(),
		btn,
	)

	border := container.NewBorder(s.Title, widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), box)
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
	s.Title = widget.NewLabelWithStyle("send binary", fyne.TextAlignCenter, fyne.TextStyle{})
	sendFunc := func() error {
		if s.serializeType != protocol.CodeNone {
			if s.BodyInput.Entry.Text == "" {
				return errors.New("body empty")
			}
		}
		if RpcClient.Client != nil {
			var body interface{}
			if s.serializeType == protocol.Json {
				body = []byte(s.BodyInput.Entry.Text)
			} else if s.serializeType == protocol.String {
				body = s.BodyInput.Entry.Text
			} else if s.serializeType == protocol.CodeNone {
				body = nil
			}
			metaData := make(map[string]string)
			if s.Metadata.Entry.Text != "" {
				err := json.Unmarshal([]byte(s.Metadata.Entry.Text), &metaData)
				if err != nil {
					return fmt.Errorf("send json.Unmarshal:err:%s", err.Error())
				}
			}

			return RpcClient.Client.Call(s.Path.Entry.Text, s.Method.Entry.Text, metaData, s.serializeType, body)
		}
		return errors.New("rpc client nil")
	}
	s.SendBtn = widget.NewButton("send", func() {

		err := sendFunc()
		if err != nil {
			GlobalText.msgChan <- err.Error()
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
		} else if tmp == "none" {
			s.serializeType = protocol.CodeNone
		}
	})
	tmpSelect.SetSelectedIndex(1)

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
	s.SendInterval = &Input{
		Layout:    layout.NewFormLayout(),
		Label:     widget.NewLabel("interval:"),
		Entry:     widget.NewEntry(),
		Container: nil,
	}
	s.SendInterval.Entry.SetPlaceHolder("ms")
	s.SendInterval.Entry.SetText("200")
	s.SendInterval.Join().Hide()
	done := make(chan struct{})
	s.Checkbox = widget.NewCheck("auto send", func(checked bool) {
		checkedFalseF := func() {
			s.SendInterval.Join().Hide()
			s.SendBtn.Show()
			done <- struct{}{}
		}
		if checked {
			s.SendInterval.Join().Show()
			s.SendBtn.Hide()

			go func() {
				for {
					select {
					case <-done:
						return
					default:
						err := sendFunc()
						if err != nil {
							GlobalText.msgChan <- err.Error()
						}
						intervalInt, err := strconv.Atoi(s.SendInterval.Entry.Text)
						if err != nil {
							GlobalText.msgChan <- err.Error()
						}
						if intervalInt < 1 {
							intervalInt = 1
						}
						//TODO 为什么这里不能直接退出Wed Dec 27 17:34:32 CST 202
						interval := time.Duration(intervalInt)
						time.Sleep(time.Millisecond * interval)
					}
				}
			}()
		} else {
			checkedFalseF()
		}

	})
	return s.Join()
}
