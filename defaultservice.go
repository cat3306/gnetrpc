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
	last, ok := ctx.Conn.GetProperty(lastHeartbeatKey)
	now := time.Now().UnixMilli()
	if ok {
		milli := last.(int64)
		rpclog.Info(ctx.Conn.Id(), ";", now-milli)
	}
	// if ok {
	// 	milli := last.(int64)
	// 	rpclog.Info(ctx.Conn.Id(), ";", now-milli)
	// 	if now-milli < maxHeartbeatFrequency.Milliseconds() {

	// 		msg := fmt.Sprintf("id:%s the heart interval is %dms should be greater than %dms closed it", ctx.Conn.Id(), now-milli, maxHeartbeatFrequency/time.Millisecond)
	// 		rpclog.Errorf(msg)
	// 		ctx.Conn.Close("heart interval")
	// 		os.Exit(0)
	// 		return nil
	// 	}
	// }
	ctx.Conn.SetProperty(lastHeartbeatKey, now)
	if b.debug {
		rpclog.Info(*args)
	}
	//rpclog.Info(*args, ctx.Metadata)
	*reply = *args
	return CallSelf()
}

func (b *BuiltinService) TestRpc(ctx *protocol.Context, args *string, reply *string) *CallMode {
	rpclog.Infof("testrpc:%s", *args)
	*reply = `\(^o^)/~`

	return CallSelf()
}
