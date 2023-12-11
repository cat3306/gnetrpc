package main

import (
	"fmt"
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"net"
	"time"
)

type CallReq struct {
	A int `json:"a"`
	B int `json:"b"`
}

func recve(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		fmt.Println(n)
	}
}
func main() {
	client, err := gnetrpc.NewClient()
	if err != nil {
		return
	}
	conn, err := client.Dial("tcp", ":7898")
	if err != nil {
		panic(err)
		return
	}
	client.Run()
	ctx := protocol.Context{
		ServicePath:   "Builtin",
		ServiceMethod: "Heartbeat",
		SerializeType: uint16(protocol.String),
		Seq:           123,
	}
	buffer := protocol.Encode(&ctx, "💓")
	for {
		conn.Write(buffer.Bytes())
		cx := <-client.CtxChan()
		fmt.Println(cx.ServicePath, cx.Payload)
		time.Sleep(time.Millisecond * 10)
	}
}
