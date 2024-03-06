package gnetrpc

import (
	"time"

	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/valyala/bytebufferpool"
)

func NewClient(address string, network string, options ...OptionFn) *Client {

	c := &Client{
		option: new(serverOption),
		pluginContainer: &pluginContainer{
			plugins: map[PluginType][]Plugin{},
		},
		address:    address,
		network:    network,
		handlerSet: NewHandlerSet(),
	}
	for _, op := range options {
		op(c.option)
	}
	if c.option.mainGoroutineChannelCap == 0 {
		c.option.mainGoroutineChannelCap = 1024
	}
	c.asyncMode = c.option.clientAsyncMode
	c.mainCtxChan = make(chan *protocol.Context, c.option.mainGoroutineChannelCap)
	cli, err := gnet.NewClient(c, gnet.WithOptions(c.option.gnetOptions))
	if err != nil {
		panic(err)
	}
	c.gCli = cli
	c.gPool = goroutine.Default()
	return c
}

type Client struct {
	option *serverOption
	*gnet.BuiltinEventEngine
	gCli            *gnet.Client
	mainCtxChan     chan *protocol.Context
	pluginContainer *pluginContainer
	handlerSet      *HandlerSet
	address         string
	network         string
	asyncMode       bool
	conn            gnet.Conn
	gPool           *goroutine.Pool
}

func (c *Client) MainGoroutine() {

	for ctx := range c.mainCtxChan {
		c.process(ctx)
	}

}
func (c *Client) process(ctx *protocol.Context) {
	err := c.handlerSet.Call(ctx, c.gPool)
	if err != nil {
		rpclog.Errorf("process err:%s,path:%s,method:%s", err.Error(), ctx.ServicePath, ctx.ServiceMethod)
	}
}
func (c *Client) Close(msg string) {
	c.conn.Close(msg)

}
func (c *Client) Run() (*Client, error) {
	_, err := c.dial()
	if err != nil {
		return nil, err
	}
	if c.asyncMode {
		go c.MainGoroutine()
	}
	_ = c.gCli.Start()
	return c, nil
}
func (c *Client) dial() (*Client, error) {
	conn, err := c.gCli.Dial(c.network, c.address)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return c, err
}
func (c *Client) OnBoot(e gnet.Engine) (action gnet.Action) {
	return
}
func (c *Client) OnShutdown(e gnet.Engine) {
}
func (c *Client) OnOpen(conn gnet.Conn) ([]byte, gnet.Action) {
	rpclog.Info("client connected,id:", conn.Fd())
	c.pluginContainer.DoDo(PluginTypeOnOpen, c.conn)
	return nil, gnet.None
}

func (c *Client) OnClose(conn gnet.Conn, err error) gnet.Action {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	c.pluginContainer.DoDo(PluginTypeOnClose, c.conn)
	rpclog.Warnf("conn close err:%s", errMsg)
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
	c.mainCtxChan <- ctx
	return gnet.None
}
func (c *Client) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
func (c *Client) Register(is ...IService) *Client {
	for _, v := range is {
		c.handlerSet.Register(v, true)
	}
	return c
}
func (c *Client) AddPlugin(ps ...Plugin) {
	for _, p := range ps {
		c.pluginContainer.Add(p.Type(), p)
	}
}
func (c *Client) Call(servicePath string, serviceMethod string, metadata map[string]string, sType protocol.SerializeType, v interface{}) error {
	ctx := protocol.Context{
		ServicePath:   servicePath,
		ServiceMethod: serviceMethod,
		Metadata:      metadata,
		H: &protocol.Header{
			MagicNumber:   protocol.MagicNumber,
			Version:       protocol.Version,
			SerializeType: uint8(sType),
		},
	}
	buffer := protocol.Encode(&ctx, v)
	defer bytebufferpool.Put(buffer)
	_, err := c.conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}
