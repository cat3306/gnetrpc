package service

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
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
	roomMgr.ConnOnClose(ctx.Conn)
}
