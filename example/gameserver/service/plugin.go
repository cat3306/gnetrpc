package service

import (
	"errors"
	"sync"

	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/gnet/v2"
)

type ConnClosePlugin struct {
}

func (c *ConnClosePlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnClose
}
func (c *ConnClosePlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return c
}
func (c *ConnClosePlugin) OnDo(args ...interface{}) error {
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
		return errors.New("args not found ctxChan")
	}
	c.doConnClose(ctxChan, conn)
	return nil
}
func (c *ConnClosePlugin) doConnClose(ctxChan chan *protocol.Context, conn gnet.Conn) {
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

type ConnOpenPlugin struct {
}

func (c *ConnOpenPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnOpen
}
func (c *ConnOpenPlugin) Init(args ...interface{}) gnetrpc.Plugin {
	return c
}
func (c *ConnOpenPlugin) OnDo(args ...interface{}) error {
	if len(args) == 0 {
		return nil
	}
	conn := args[0].(gnet.Conn)
	conn.SetContext(&sync.Map{}) //
	return nil
}
