package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/share"
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
func (b *Builtin) Alias() string {
	return ""
}
func (b *Builtin) Heartbeat(ctx *protocol.Context) {
	rpclog.Info(ctx.Payload.String())
}

type Account struct {
	salt   string
	secret string
}

func (a *Account) Alias() string {
	return ""
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
func (a *Account) Register(ctx *protocol.Context) {
	rpclog.Info(ctx.Payload.String())
}
func (a *Account) EmailCode(ctx *protocol.Context) {
	rpclog.Info(ctx.Payload.String())
}
func multiClient(num int) {
	for i := 0; i < num; i++ {
		go singleClient()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
	fmt.Println("done")
	select {}
}
func singleClient() {
	client, err := gnetrpc.NewClient("127.0.0.1:7898", "tcp", gnetrpc.WithClientAsyncMode()).
		Register(
			new(Builtin),
			new(Account),
		).Run()
	if err != nil {
		panic(err)
	}

	HeartBeat(client)
}
func main() {
	multiClient(1)
	//singleClient()
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

type RegisterReq struct {
	Email string `json:"Email"`
	Pwd   string `json:"Pwd"`
	Nick  string `json:"Nick"`
	Code  string `json:"Code"`
}

func Register(client *gnetrpc.Client) {
	req := RegisterReq{
		Email: "2696584197@qq.com",
		Pwd:   "123",
		Nick:  "cat101",
		Code:  "091533",
	}
	err := client.Call("Account", "Register", nil, protocol.Json, req)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10)
}
func Logout(client *gnetrpc.Client) {
	err := client.Call("Account", "Logout", nil, protocol.CodeNone, nil)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10)
}

type EmailCodeReq struct {
	Email string `json:"email"`
}

func EmailCode(client *gnetrpc.Client) {
	req := EmailCodeReq{
		Email: "2696584197@qq.com",
	}
	err := client.Call("Account", "EmailCode", nil, protocol.Json, req)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10)
}
func HeartBeat(client *gnetrpc.Client) {

	for i := 0; i < 10; i++ {
		err := client.Call("Builtin", "Heartbeat", map[string]string{
			share.AuthKey: "鸳鸯擦，鸳鸯体，你爱我，我爱你",
		}, protocol.String, "💓:"+strconv.Itoa(i))
		if err != nil {
			break
		}
	}
}
