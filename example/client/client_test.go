package main

import (
	"github.com/cat3306/gnetrpc/protocol"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":7898")
	if err != nil {
		t.Fatalf(err.Error())
	}
	req := CallReq{
		A: 10,
		B: 12,
	}
	ctx := protocol.Context{
		ServicePath:   "Arith",
		ServiceMethod: "Add",
		SerializeType: uint16(protocol.Json),
		Seq:           123,
	}
	buffer := protocol.Encode(&ctx, &req)
	conn.Write(buffer.Bytes())
	select {}
}
