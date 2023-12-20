package main

import (
	"fmt"
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

type Account struct {
	salt   string
	secret string
}

func (a *Account) Init(v ...interface{}) gnetrpc.IService {
	return a
}

func (a *Account) Login(ctx *protocol.Context) {
	rpclog.Info(ctx.Payload.String())
}
func (a *Account) Logout(ctx *protocol.Context) {
	rpclog.Info(ctx.Payload.String())
}
func main() {
	client, err := gnetrpc.NewClient("127.0.0.1:7898", "tcp", gnetrpc.WithClientAsyncMode()).
		Register(
			new(Builtin),
			new(Account),
		).Run()
	if err != nil {
		panic(err)
	}
	Login(client)
	//HeartBeat(client)
}

type LoginReq struct {
	Email string `json:"Email"`
	Pwd   string `json:"Pwd"`
}

func Login(client *gnetrpc.Client) {
	req := LoginReq{
		Email: "1273014435@qq.com",
		Pwd:   "123",
	}
	err := client.Call("Account", "Login", nil, protocol.Json, req)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10)
	Logout(client)
}
func Logout(client *gnetrpc.Client) {
	err := client.Call("Account", "Logout", nil, protocol.CodeNone, nil)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10)
}
func HeartBeat(client *gnetrpc.Client) {
	for {
		err := client.Call("Builtin", "Heartbeat", map[string]string{
			share.AuthKey: "鸳鸯擦，鸳鸯体，你爱我，我爱你",
		}, protocol.String, "💓")
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
