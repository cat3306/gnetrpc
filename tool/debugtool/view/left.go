package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/rpclog"
)

var (
	NetWorkSelect = []string{gnetrpc.TcpNetwork, gnetrpc.UdpNetwork, gnetrpc.UnixNetwork}
	RpcClient     ClientInfo
)

type ClientInfo struct {
	Ip      string `json:"Ip"`
	Port    string `json:"port"`
	Network string `json:"network"`
	Client  *gnetrpc.Client
}
type UpperLeftView struct {
	Title      *widget.Label
	IpInput    *Input
	PortInput  *Input
	Container  *fyne.Container
	NetWork    *Select
	ConnBtn    *widget.Button
	DisConnBtn *widget.Button
	ClientInfo ClientInfo
}

func (u *UpperLeftView) Join() *fyne.Container {
	if u.Container == nil {

		//btnBox := container.NewBorder(nil, nil, u.ConnBtn, u.DisConnBtn, nil)
		buttonsContainer := container.New(layout.NewHBoxLayout(),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(155, 30)), u.ConnBtn),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(155, 30)), u.DisConnBtn),
		)
		box := container.NewVBox(u.Title,
			u.IpInput.Join(),
			u.PortInput.Join(),
			u.NetWork.Join(),
			buttonsContainer,
		)
		border := container.NewBorder(widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), box)
		u.Container = border
	}
	return u.Container
}
func (u *UpperLeftView) Check() bool {
	if u.IpInput.Entry.Text == "" || u.PortInput.Entry.Text == "" {
		return false
	}
	return true
}
func (u *UpperLeftView) DisConnect() {
	defer func() {
		u.ConnBtn.Enable()
		u.DisConnBtn.Disable()
	}()
	RpcClient.Client.Close("closed by manual")
}
func (u *UpperLeftView) Connect() {

	if !u.Check() {
		return
	}
	u.ClientInfo.Ip = u.IpInput.Entry.Text
	u.ClientInfo.Port = u.PortInput.Entry.Text
	addr := fmt.Sprintf("%s:%s", u.IpInput.Entry.Text, u.PortInput.Entry.Text)
	client, err := gnetrpc.
		NewClient(addr, u.ClientInfo.Network).
		Run()
	rpclog.SetLogger(&Log{})
	if err != nil {
		GlobalText.Chan() <- err.Error()
		return
	}
	RpcClient.Client = client
	cp := &ClosePlugin{
		ConnectBtn:    u.ConnBtn,
		DisConnectBtn: u.DisConnBtn,
	}
	op := &OpenPlugin{
		ConnectBtn:    u.ConnBtn,
		DisConnectBtn: u.DisConnBtn,
	}
	RpcClient.Client.AddPlugin(cp, op)
	go u.ReceiveDemon()
}
func (u *UpperLeftView) ReceiveDemon() {
	for ctx := range RpcClient.Client.CtxChan() {
		ReceiveText.msgChan <- ctx.Payload.String()
	}
}
func InitUpperLeftView() *UpperLeftView {
	title := widget.NewLabelWithStyle("client setting", fyne.TextAlignCenter, fyne.TextStyle{})
	ie := widget.NewEntry()
	ie.SetText("127.0.0.1")
	ip := &Input{
		Layout: layout.NewFormLayout(),
		Label:  widget.NewLabel("ip:"),
		Entry:  ie,
	}

	pe := widget.NewEntry()
	pe.SetText("7898")
	port := &Input{
		Layout: layout.NewFormLayout(),
		Label:  widget.NewLabel("port:"),
		Entry:  pe,
	}

	uv := &UpperLeftView{
		IpInput:   ip,
		PortInput: port,
		Title:     title,
	}
	tmpNetwork := widget.NewSelect(NetWorkSelect, uv.SelectNetwork)
	tmpNetwork.SetSelected(NetWorkSelect[0])
	uv.NetWork = &Select{
		Layout: layout.NewFormLayout(),
		Label:  widget.NewLabel("network:"),
		Select: tmpNetwork,
	}

	uv.ConnBtn = widget.NewButton("connect", uv.Connect)
	uv.DisConnBtn = widget.NewButton("disconnect", uv.DisConnect)
	uv.DisConnBtn.Disable()
	return uv
}
func (u *UpperLeftView) SelectNetwork(mode string) {
	u.ClientInfo.Network = mode
}
