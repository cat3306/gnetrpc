package gnetrpc

import (
	"time"

	"github.com/cat3306/gnetrpc/protocol"
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

func (b *BuiltinService) Benchmark(ctx *protocol.Context, args *protocol.BenchmarkMessage, reply *protocol.BenchmarkMessage, tag struct{}) *CallMode {
	*reply = *args
	//rpclog.Info(args.Field1)
	return CallSelf()
}
