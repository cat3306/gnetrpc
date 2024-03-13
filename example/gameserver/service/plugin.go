package service

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/gnet/v2"
)

type ClosePlugin struct {
}

func (c *ClosePlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnClose
}
func (c *ClosePlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return c
}
func (c *ClosePlugin) OnDo(args ...interface{}) interface{} {
	conn := args[0].(gnet.Conn)
	rpclog.Warnf("client close id:%s,cause:%v", util.GetConnId(conn), args[1])
	var (
		ctxChan chan *protocol.Context
		ok      bool
	)
	for _, v := range args {
		ctxChan, ok = v.(chan *protocol.Context)
		if ok {
			break
		}
	}
	if ctxChan == nil {
		rpclog.Errorf("args not found ctxChan")
		return true
	}
	c.doConnClose(ctxChan, conn)
	return true
}
func (c *ClosePlugin) doConnClose(ctxChan chan *protocol.Context, conn gnet.Conn) {
	ctx := protocol.GetCtx()
	ctx.H.SerializeType = byte(protocol.Json)
	ctx.H.Version = protocol.Version
	ctx.H.HeartBeat = 0
	ctx.H.MagicNumber = protocol.MagicNumber
	ctx.ServicePath = "GameMgr"
	ctx.ServiceMethod = "ConnOnClose"
	ctx.Conn = conn
	ctxChan <- ctx
}
