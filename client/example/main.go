package main

import (
	"github.com/cat3306/gnetrpc/protocol"
	"net"
	"time"
)

type CallReq struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	conn, err := net.Dial("tcp", ":7898")
	if err != nil {
		panic(err)
	}
	req := CallReq{
		A: 10,
		B: 12,
	}
	ctx := protocol.Context{
		ServicePath:   "Arith",
		ServiceMethod: "AsyncMakeLove",
		SerializeType: uint16(protocol.Json),
		Seq:           123,
	}
	buffer := protocol.Encode(&ctx, &req)
	for {
		conn.Write(buffer.Bytes())
		time.Sleep(time.Millisecond * 10)
	}
}
