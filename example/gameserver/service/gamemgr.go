package service

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/example/gameserver/util"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
)

type GameMgr struct {
}

func (r *GameMgr) Init(v ...interface{}) gnetrpc.IService {
	return r
}
func (r *GameMgr) Alias() string {
	return ""
}

func (r *GameMgr) ConnOnClose(ctx *protocol.Context) {
	err := roomMgr.ConnOnClose(ctx.Conn)
	if err != nil {
		rpclog.Infof("ConnOnClose,err:%s", err.Error())
	}
}
func (r *GameMgr) GlobalChat(ctx *protocol.Context, req *string, rsp *string) *gnetrpc.CallMode {
	*rsp = *req
	m := util.GetSyncMapFromConn(ctx.Conn)
	m.Store("haha", "gaigai")
	return gnetrpc.CallBroadcast()
}
