package main

import (
	"context"
	"fmt"
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
)

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	fmt.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func (t *Arith) Add(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	fmt.Printf("call: %d + %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func (t *Arith) Say(ctx context.Context, args *string, reply *string) error {
	*reply = "hello " + *args
	return nil
}
func (t *Arith) MakeLove(ctx *protocol.Context) (*gnetrpc.CallMode, error) {
	return nil, nil
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
