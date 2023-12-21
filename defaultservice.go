package gnetrpc

import (
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"time"
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
func (b *BuiltinService) Heartbeat(ctx *protocol.Context, args *string, reply *string) *CallMode {
	last, ok := ctx.Conn.GetProperty(lastHeartbeatKey)
	if ok {
		now := time.Now().UnixMilli()
		milli := last.(int64)
		if now-milli < maxHeartbeatFrequency.Milliseconds() {
			msg := fmt.Sprintf("the heart interval should be greater than %dms closed it", maxHeartbeatFrequency/time.Millisecond)
			rpclog.Errorf(msg)
			ctx.Conn.Close("heart interval")
			return nil
		}
	}
	ctx.Conn.SetProperty(lastHeartbeatKey, time.Now().UnixMilli())
	rpclog.Info(*args)
	*reply = "❤️"
	return CallSelf()
}
