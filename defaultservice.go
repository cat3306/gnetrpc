package gnetrpc

import (
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
)

const BuiltinServiceName = "BuiltinService"

type BuiltinService struct {
}

func (b *BuiltinService) Heartbeat(ctx *protocol.Context, args *string, reply *string) (*CallMode, error) {
	rpclog.Info(*args)
	*reply = "❤️"
	return CallSelf(), nil
}

type AuthReq struct {
}

func (b *BuiltinService) Auth(ctx *protocol.Context, args *string, reply *string) (*CallMode, error) {
	rpclog.Info(*args)
	*reply = "❤️"
	return CallSelf(), nil
}
