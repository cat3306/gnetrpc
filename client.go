package gnetrpc

import (
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/panjf2000/gnet/v2"
	"time"
)

func NewClient(options ...OptionFn) (*Client, error) {
	c := &Client{
		mainCtxChan: make(chan *protocol.Context, 1024),
		option:      new(serverOption),
	}
	for _, op := range options {
		op(c.option)
	}
	cli, err := gnet.NewClient(c, gnet.WithOptions(c.option.gnetOptions))
	if err != nil {
		return nil, err
	}
	c.gCli = cli
	return c, nil
}

type Client struct {
	option *serverOption
	*gnet.BuiltinEventEngine
	gCli        *gnet.Client
	mainCtxChan chan *protocol.Context
}

func (c *Client) MainGoroutine() {

	for ctx := range c.mainCtxChan {
		c.process(ctx)
	}

}
func (c *Client) process(ctx *protocol.Context) {

}
func (c *Client) Run() {
	_ = c.gCli.Start()
}
func (c *Client) Dial(network, address string) (gnet.Conn, error) {
	return c.gCli.Dial(network, address)
}
func (c *Client) OnBoot(e gnet.Engine) (action gnet.Action) {
	return
}
func (c *Client) OnShutdown(e gnet.Engine) {
}
func (c *Client) OnOpen(conn gnet.Conn) ([]byte, gnet.Action) {
	rpclog.Info("conn", conn.Fd())
	return nil, gnet.None
}

func (c *Client) OnClose(conn gnet.Conn, err error) gnet.Action {

	return gnet.None
}
func (c *Client) CtxChan() <-chan *protocol.Context {
	return c.mainCtxChan
}
func (c *Client) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	ctx, err := protocol.Decode(conn)
	if err != nil {
		return
	}
	//rpclog.Infof("%s,%s,%s", ctx.ServicePath, ctx.ServiceMethod, string(ctx.Payload.B))
	c.mainCtxChan <- ctx
	return gnet.None
}
func (c *Client) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
