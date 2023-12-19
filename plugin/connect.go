package plugin

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/panjf2000/gnet/v2"
)

type ConnectPlugin struct {
}

func (c *ConnectPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnOpen
}
func (c *ConnectPlugin) OnDo(args interface{}) interface{} {
	conn := args.(gnet.Conn)
	rpclog.Infof("client connect cid:%s", conn.Id())
	return true
}
