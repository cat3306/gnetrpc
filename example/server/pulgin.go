package main

import (
	"fmt"
	"github.com/cat3306/gnetrpc"
)

type BlackListPlugin struct {
}

func (b *BlackListPlugin) OnDo(v interface{}) interface{} {
	fmt.Println("test BlackListPlugin")
	return nil
}
func (b *BlackListPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnOpen
}
