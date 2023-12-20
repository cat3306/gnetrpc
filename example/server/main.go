package main

import (
	"errors"
	"fmt"
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/plugin"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
)

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Reply struct {
	C int
}

type Arith struct {
	set map[string]string
}

func (t *Arith) Init(v ...interface{}) gnetrpc.IService {
	return t
}
func (t *Arith) Mul(ctx *protocol.Context, args *Args, reply *Reply) (*gnetrpc.CallMode, error) {
	reply.C = args.A * args.B
	fmt.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return gnetrpc.CallSelf(), nil
}

func (t *Arith) Add(ctx *protocol.Context, args *Args, reply *Reply) (*gnetrpc.CallMode, error) {
	reply.C = args.A + args.B
	fmt.Printf("call: %d + %d = %d\n", args.A, args.B, reply.C)
	return nil, nil
}

func (t *Arith) Say(ctx *protocol.Context, args *string, reply *string, tag struct{}) (*gnetrpc.CallMode, error) {
	*reply = "hello " + *args
	rpclog.Info(*reply)
	return nil, nil
}
func (t *Arith) MakeLove(ctx *protocol.Context) {
	rpclog.Info(ctx.MsgSeq)
}
func (t *Arith) AsyncMakeLove(ctx *protocol.Context, tag struct{}) {
	rpclog.Info("AsyncMakeLove")
}
func main() {
	s := gnetrpc.NewServer(
		gnetrpc.WithMulticore(true),
		gnetrpc.WithPrintRegisteredMethod(),
		gnetrpc.WithDefaultService(),
		gnetrpc.WithMainGoroutineChannelCap(10000),
	)
	s.UseAuthFunc(func(ctx *protocol.Context, token string) error {
		if token != "鸳鸯擦，鸳鸯体，你爱我，我爱你" {
			return errors.New("你不爱我 !")
		}
		return nil
	})
	s.AddPlugin(
		new(plugin.ConnectPlugin),
		new(plugin.ClosePlugin),
		new(plugin.BlacklistPlugin).Add("211.137.99.189"),
	)
	s.Register(new(Arith))
	err := s.Run(gnetrpc.TcpNetwork, ":7898")
	fmt.Println(err)
}
