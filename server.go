package gnetrpc

import (
	"errors"
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"runtime/debug"
	"sync"
	"time"
)

type serverOption struct {
	printMethod    bool
	defaultService bool
	gnetOptions    gnet.Options
	antOption      ants.Options
}
type Server struct {
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	serviceMapMu sync.RWMutex
	serviceSet   *ServiceSet
	handlerSet   *HandlerSet
	option       *serverOption
	gPool        *ants.Pool
	mainCtxChan  chan *protocol.Context
	connMatrix   *connMatrix
}

func (s *Server) MainGoroutine() {

	for ctx := range s.mainCtxChan {
		s.process(ctx)
	}

}

func (s *Server) process(ctx *protocol.Context) {
	defer func() {
		if r := recover(); r != nil {
			msg := debug.Stack()
			err := fmt.Errorf("[server call internal error] service: %s, method: %s, stack: %s,err:%s", ctx.ServicePath, ctx.ServiceMethod, util.BytesToString(msg), r)
			rpclog.Error(err)
		}
	}()

	servicePath := ctx.ServicePath
	method := ctx.ServiceMethod
	err := s.handlerSet.ExecuteHandler(ctx, s.gPool)
	if err == nil {
		return
	}
	if err != nil && !errors.Is(err, NotFoundMethod) {
		rpclog.Errorf("process err:%s,service:%s, method:%s", err.Error(), servicePath, method)
		return
	}
	err = s.serviceSet.Call(ctx, s)
	if err != nil {
		rpclog.Errorf("process err:%s,service:%s, method:%s", err.Error(), servicePath, method)
	}

}
func (s *Server) OnBoot(engine gnet.Engine) (action gnet.Action) {
	return
}

func (s *Server) OnShutdown(engine gnet.Engine) {
}

func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	rpclog.Infof("client connect cid:%d", c.Fd())
	s.connMatrix.Add(c)
	return
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	s.connMatrix.Remove(c.Fd())
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	rpclog.Infof("cid:%d close,reason:%s", c.Fd(), reason)
	return
}

func (s *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx, err := protocol.Decode(c)
	if err != nil {
		rpclog.Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	s.mainCtxChan <- ctx
	return
}
func (s *Server) Register(v IService, name ...string) {
	s.serviceSet.Register(v, s.option.printMethod, name...)
	s.registerRouter(v, name...)
}

// registerRouter func(ctx *protocol.Context) or func(ctx *protocol.Context, tag struct{}) should by registered
func (s *Server) registerRouter(v IService, name ...string) {
	s.handlerSet.Register(v, s.option.printMethod, name...)
}

func (s *Server) OnTick() (delay time.Duration, action gnet.Action) {
	return
}

func (s *Server) Run(netWork string, addr string) error {
	go s.MainGoroutine()
	return gnet.Run(s, netWork+"://"+addr, gnet.WithOptions(s.option.gnetOptions))
}
func NewServer(options ...OptionFn) *Server {
	s := &Server{
		//Plugins:    &pluginContainer{},
		//options:    make(map[string]interface{}),
		//activeConn: make(map[net.Conn]struct{}),
		//doneChan:   make(chan struct{}),
		serviceSet: NewServiceSet(),
		handlerSet: NewHandlerSet(),
		//router:     make(map[string]Handler),
		gPool:       goroutine.Default(),
		mainCtxChan: make(chan *protocol.Context, 1024),
		connMatrix:  newConnMatrix(),
		option:      new(serverOption),
	}

	for _, op := range options {
		op(s.option)
	}

	//if s.options["TCPKeepAlivePeriod"] == nil {
	//	s.options["TCPKeepAlivePeriod"] = 3 * time.Minute
	//}
	if s.option.defaultService {
		s.Register(new(BuiltinService), "Builtin")
	}
	return s
}
