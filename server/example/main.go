package main

import (
	"fmt"
	"github.com/cat3306/gnetrpc"
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

func (t *Arith) Mul(ctx *protocol.Context, args *Args, reply *Reply) (*gnetrpc.CallMode, error) {
	reply.C = args.A * args.B
	fmt.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil, nil
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
	rpclog.Info(ctx.Seq)
}
func (t *Arith) AsyncMakeLove(ctx *protocol.Context, tag struct{}) {
	rpclog.Info("AsyncMakeLove")
}
func main() {
	s := gnetrpc.NewServer(
		gnetrpc.WithMulticore(true),
		gnetrpc.WithPrintRegisteredMethod(),
	)
	s.Register(new(Arith))
	s.RegisterRouter(new(Arith))
	err := s.Run(gnetrpc.TcpNetwork, ":7898")
	fmt.Println(err)
}
