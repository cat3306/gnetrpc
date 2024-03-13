package plugin

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/gnet/v2"
)

type ConnectPlugin struct {
}

func (c *ConnectPlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return c
}
func (c *ConnectPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnOpen
}
func (c *ConnectPlugin) OnDo(args ...interface{}) error {
	if len(args) == 0 {
		return nil
	}
	conn := args[0].(gnet.Conn)
	rpclog.Infof("client connect id:%s", util.GetConnId(conn))
	return nil
}
