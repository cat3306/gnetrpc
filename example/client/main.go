package main

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/share"
	"time"
)

type CallReq struct {
	A int `json:"a"`
	B int `json:"b"`
}
type Builtin struct {
}

func (b *Builtin) Init(v ...interface{}) gnetrpc.IService {
	return b
}
func (b *Builtin) Heartbeat(ctx *protocol.Context) {
	rpclog.Info(ctx.Payload.String())
}
func main() {
	client, err := gnetrpc.NewClient("183.232.230.25:7898", "tcp", gnetrpc.WithClientAsyncMode()).
		Register(new(Builtin)).Run()
	if err != nil {
		panic(err)
	}
	HeartBeat(client)
}

func HeartBeat(client *gnetrpc.Client) {
	for {
		err := client.Call("Builtin", "Heartbeat", map[string]string{
			share.AuthKey: "鸳鸯擦，鸳鸯体，你爱我，我爱你",
		}, protocol.String, "💓")
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}
