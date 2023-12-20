package plugin

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/panjf2000/gnet/v2"
)

type ClosePlugin struct {
}

func (c *ClosePlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnClose
}
func (c *ClosePlugin) OnDo(args ...interface{}) interface{} {
	conn := args[0].(gnet.Conn)
	rpclog.Warnf("client close id:%s,cause:%v", conn.Id(), args[1])
	return true
}
