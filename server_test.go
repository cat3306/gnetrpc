package gnetrpc

import (
	"context"
	"fmt"
	"testing"
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
func TestServiceSet_Register(t *testing.T) {
	s := NewServiceSet()
	s.Register(new(Arith), false)

}

func TestServer(t *testing.T) {
	s := NewServer(WithMulticore(true), WithPrintRegisteredMethod())
	s.Register(new(Arith))
	err := s.Run(TcpNetwork, ":7898")
	t.Fatalf(err.Error())
}
