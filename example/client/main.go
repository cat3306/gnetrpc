package main

import (
	"fmt"
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/share"
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
		Metadata: map[string]string{
			"abc":         "123",
			"name":        "joker",
			share.AuthKey: "鸳鸯擦，鸳鸯体，你爱我，我爱你",
		},
		H: &protocol.Header{
			MagicNumber:   protocol.MagicNumber,
			Version:       protocol.Version,
			HeartBeat:     0,
			SerializeType: uint8(protocol.String),
		},
		MsgSeq: 123,
	}
	buffer := protocol.Encode(&ctx, "💓")
	for {
		conn.Write(buffer.Bytes())
		time.Sleep(time.Millisecond * 1000)
	}
}
