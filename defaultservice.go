package gnetrpc

import (
	"time"

	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
)

const (
	BuiltinServiceName    = "Builtin"
	maxHeartbeatFrequency = time.Millisecond * 100
	lastHeartbeatKey      = "BeatTime"
)

type BuiltinService struct {
	debug bool
}

func (b *BuiltinService) Init(v ...interface{}) IService {
	return b
}
func (b *BuiltinService) Alias() string {
	return BuiltinServiceName
}
func (b *BuiltinService) Heartbeat(ctx *protocol.Context, args *string, reply *string, tag struct{}) *CallMode {
	*reply = *args
	return CallSelf()
}

func (b *BuiltinService) TestRpc(ctx *protocol.Context, args *string, reply *string) *CallMode {
	rpclog.Infof("testrpc:%s", *args)
	*reply = `\(^o^)/~`

	return CallSelf()
}
