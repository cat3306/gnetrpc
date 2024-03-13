package service

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/rpclog"
)

type ShutdownPlugin struct {
}

func (s *ShutdownPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnShutdown
}
func (s *ShutdownPlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return s
}
func (s *ShutdownPlugin) OnDo(args ...interface{}) error {
	rpclog.Infof("game shutdown")
	return nil
}
