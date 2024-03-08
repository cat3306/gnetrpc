package gnetrpc

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cat3306/gnetrpc/common"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/share"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/valyala/bytebufferpool"
)

type serverOption struct {
	printMethod             bool
	defaultService          bool
	mainGoroutineChannelCap int
	clientAsyncMode         bool //client
	gnetOptions             gnet.Options
	antOption               ants.Options
}
type Server struct {
	gnet.BuiltinEventEngine
	network         string
	addr            string
	eng             gnet.Engine
	serviceSet      *ServiceSet
	option          *serverOption
	gPool           *ants.Pool
	mainCtxChan     chan *protocol.Context
	connMatrix      *common.ConnMatrix
	authFunc        func(ctx *protocol.Context, token string) error
	hotHandlerNum   int32
	pluginContainer *pluginContainer
	signalChan      chan os.Signal
	shutdownCtx     context.Context
	shutdonwFunc    context.CancelFunc
}

func (s *Server) UseAuthFunc(f func(ctx *protocol.Context, token string) error) {
	s.authFunc = f
}

func (s *Server) mainGoroutine() {
	for {
		select {
		case ctx := <-s.mainCtxChan:
			servicePath := ctx.ServicePath
			method := ctx.ServiceMethod
			err := s.serviceSet.SyncCall(ctx)
			if err != nil {
				rpclog.Errorf("syncCall err:%s,service:%s, method:%s", err.Error(), servicePath, method)
			}
		case <-s.shutdownCtx.Done():

			rpclog.Infof("mainGoroutine shutdown")
			return
		}
	}
}

func (s *Server) process(ctx *protocol.Context) {

	servicePath := ctx.ServicePath
	method := ctx.ServiceMethod
	err := s.serviceSet.AsyncCall(ctx)
	if err != nil {
		if errors.Is(err, NotFoundMethod) {
			s.mainCtxChan <- ctx
			return
		}
		rpclog.Errorf("process err:%s,service:%s, method:%s", err.Error(), servicePath, method)
	}
}
func (s *Server) OnBoot(engine gnet.Engine) (action gnet.Action) {
	s.eng = engine
	go s.mainGoroutine()
	go s.notifySign()
	rpclog.Infof("gnetrpc start %s server on %s", s.network, s.addr)
	s.pluginContainer.DoDo(PluginTypeOnBoot, nil)
	return
}

func (s *Server) OnShutdown(engine gnet.Engine) {
	rpclog.Infof("gnetrpc shutdown")
	s.pluginContainer.DoDo(PluginTypeOnShutdown, nil)
	s.connMatrix.RemoveAll("gnetrpc shutdown")
}

func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetId(util.GenConnId(c.Fd()))
	s.connMatrix.Add(c)
	plugins := s.pluginContainer.Plugins(PluginTypeOnOpen)
	for _, v := range plugins {
		ok := v.OnDo(c).(bool)
		if !ok {
			c.Close("plugin check failed")
			return
		}
	}
	return
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	s.connMatrix.Remove(c.Id())
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	s.pluginContainer.DoDo(PluginTypeOnClose, c, reason, s.mainCtxChan)
	return
}

func (s *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	for {
		ctx, err := protocol.Decode(c)
		if err != nil {
			//rpclog.Warnf("OnTraffic err:%s", err.Error())
			break
		}
		s.pluginContainer.DoDo(PluginTypeOnTraffic, nil)
		if s.authFunc != nil {
			token := ctx.Metadata[share.AuthKey]
			err = s.authFunc(ctx, token)
			if err != nil {
				rpclog.Errorf("auth failed for conn %s: %s", c.RemoteAddr().String(), err.Error())
				err = c.Close("auth failed")
				if err != nil {
					rpclog.Errorf("conn close err:%s,%s", err.Error(), c.RemoteAddr().String())
				}
				return
			}
		}
		ctx.GPool = s.gPool
		ctx.ConnMatrix = s.connMatrix
		s.process(ctx)
	}
	return gnet.None
}
func (s *Server) Register(is ...IService) {
	for _, v := range is {
		s.serviceSet.Register(v, s.option.printMethod)
	}
}

func (s *Server) OnTick() (delay time.Duration, action gnet.Action) {
	return
}

func (s *Server) notifySign() {
	signal.Notify(s.signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-s.signalChan
		rpclog.Infof("gnetrpc notify signal:%s", sig.String())
		s.shutdonwFunc()
		err := s.eng.Stop(context.Background())
		if err != nil {
			rpclog.Errorf("notifySign err:%s", err.Error())
		}
	}()
}
func (s *Server) Run(netWork string, addr string) error {
	s.addr = addr
	s.network = netWork
	return gnet.Run(s, netWork+"://"+addr, gnet.WithOptions(s.option.gnetOptions))
}
func (s *Server) AddPlugin(ps ...Plugin) {
	for _, p := range ps {
		s.pluginContainer.Add(p.Type(), p)
	}
}
func (s *Server) SendMessage(conn gnet.Conn, path, method string, metadata map[string]string, body []byte) {
	ctx := protocol.GetCtx()
	ctx.H.MagicNumber = protocol.MagicNumber
	ctx.H.Version = protocol.Version
	ctx.H.SerializeType = byte(protocol.CodeNone)
	ctx.Conn = conn
	ctx.ServicePath = path
	ctx.ServiceMethod = method
	ctx.Metadata = metadata
	buffer := protocol.Encode(ctx, body)
	defer func() {
		bytebufferpool.Put(buffer)
	}()
	_, err := conn.Write(buffer.Bytes())
	if err != nil {
		rpclog.Errorf("SendMessage err:%s", err.Error())
	}
}
func NewServer(options ...OptionFn) *Server {
	s := &Server{
		gPool:      goroutine.Default(),
		connMatrix: common.NewConnMatrix(true),
		option:     new(serverOption),
		pluginContainer: &pluginContainer{
			plugins: map[PluginType][]Plugin{},
		},
		signalChan: make(chan os.Signal, 1),
	}
	s.shutdownCtx, s.shutdonwFunc = context.WithCancel(context.Background())
	s.serviceSet = NewServiceSet(s.gPool)
	for _, op := range options {
		op(s.option)
	}
	if s.option.mainGoroutineChannelCap == 0 {
		s.option.mainGoroutineChannelCap = 1024
	}
	s.mainCtxChan = make(chan *protocol.Context, s.option.mainGoroutineChannelCap)
	if s.option.defaultService {
		s.Register(new(BuiltinService))
	}
	return s
}
