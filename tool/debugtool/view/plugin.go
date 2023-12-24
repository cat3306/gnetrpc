package view

import (
	"fyne.io/fyne/v2/widget"
	"github.com/cat3306/gnetrpc"
)

type ClosePlugin struct {
	ConnectBtn    *widget.Button
	DisConnectBtn *widget.Button
}

func (c *ClosePlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnClose
}
func (c *ClosePlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return c
}
func (c *ClosePlugin) OnDo(args ...interface{}) interface{} {
	c.ConnectBtn.Enable()
	c.DisConnectBtn.Disable()
	return true
}

type OpenPlugin struct {
	ConnectBtn    *widget.Button
	DisConnectBtn *widget.Button
}

func (o *OpenPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnOpen
}
func (o *OpenPlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return o
}
func (o *OpenPlugin) OnDo(args ...interface{}) interface{} {
	o.ConnectBtn.Disable()
	o.DisConnectBtn.Enable()
	return true
}
