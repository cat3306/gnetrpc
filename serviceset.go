package gnetrpc

import (
	"errors"
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"reflect"
)

type ServiceSet struct {
	set    map[string]*Service
	server *Server
}

func NewServiceSet(server *Server) *ServiceSet {
	return &ServiceSet{
		set:    map[string]*Service{},
		server: server,
	}
}
func (s *ServiceSet) GetService(sp string) (bool, *Service) {
	v, ok := s.set[sp]
	return ok, v
}
func (s *ServiceSet) suitableMethods(typ reflect.Type, reportErr bool) (map[string]*methodType, map[string]*methodType) {
	methods := make(map[string]*methodType)
	asyncMethods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver, *protocol.Context, *args, *reply or struct
		if !(mtype.NumIn() == 4 || mtype.NumIn() == 5) {
			continue
		}
		// First arg must be protocol.Context
		ctxType := mtype.In(1)
		if !ctxType.AssignableTo(typeOfContext) {
			continue
		}

		// Second arg need not be a pointer.
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			continue
		}
		// Third arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		if mtype.NumIn() == 5 {
			structType := mtype.In(4)
			if structType != typeOfSturct {
				continue
			}
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			continue
		}

		// Method needs one out.
		if mtype.NumOut() != 2 {
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfCallMode {
			continue
		}
		if returnType := mtype.Out(1); returnType != typeOfError {

			continue
		}
		if mtype.NumIn() == 4 {
			methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
		}
		if mtype.NumIn() == 5 {
			asyncMethods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
		}
		// init pool for reflect.Type of args and reply
		reflectTypePools.Init(argType)
		reflectTypePools.Init(replyType)
	}
	return methods, asyncMethods
}
func (s *ServiceSet) Register(v IService, isPrint bool, name ...string) {
	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	sName := ""
	if len(name) != 0 {
		sName = name[0]
	} else {
		sName = reflect.Indirect(value).Type().Name()
	}
	if sName == "" {
		errorStr := "Register: no service name for type " + typ.String()
		panic(errorStr)
	}
	methodSet, asyncMethods := s.suitableMethods(typ, true)
	tmpService := &Service{
		name:        sName,
		value:       value,
		typ:         typ,
		method:      methodSet,
		asyncMethod: asyncMethods,
	}
	if isPrint {
		for _, m := range tmpService.method {
			rpclog.Info(fmt.Sprintf("registered [%s.%s]", tmpService.name, m.method.Name))
		}
		for _, m := range tmpService.asyncMethod {
			rpclog.Info(fmt.Sprintf("registered [%s.go_%s]", tmpService.name, m.method.Name))
		}

	}
	s.set[tmpService.name] = tmpService
}
func (s *ServiceSet) Call(ctx *protocol.Context) error {
	defer func() {
		protocol.PutCtx(ctx)
	}()
	servicePath := ctx.ServicePath
	method := ctx.ServiceMethod
	tmpService := s.set[servicePath]
	if tmpService == nil {
		return NotFoundService
	}
	var isAsync bool
	mType := tmpService.method[method]
	if mType == nil {
		isAsync = true
		mType = tmpService.asyncMethod[method]
	}
	if mType == nil {
		return NotFoundMethod
	}
	codec := protocol.GetCodec(protocol.SerializeType(ctx.H.SerializeType))
	if codec == nil {
		return errors.New("invalid serialize type")
	}
	f := func() error {
		replyv := reflectTypePools.Get(mType.ReplyType)
		argv := reflectTypePools.Get(mType.ArgType)
		defer func() {
			reflectTypePools.Put(mType.ArgType, argv)
			reflectTypePools.Put(mType.ReplyType, replyv)
		}()
		err := codec.Unmarshal(ctx.Payload.Bytes(), argv)
		if err != nil {
			return err
		}
		callModel, err := tmpService.call(ctx, mType, reflect.ValueOf(argv), reflect.ValueOf(replyv), isAsync)
		if err != nil {
			return err
		}
		if callModel == nil {
			return errors.New("call mode nil")
		}

		switch callModel.Call {
		case None:
			return nil
		case Self:
			buffer := protocol.Encode(ctx, replyv)
			s.server.connMatrix.SendToOne(buffer, ctx.Conn.Id())
			return nil
		case Broadcast:
			buffer := protocol.Encode(ctx, replyv)
			s.server.connMatrix.Broadcast(buffer)
		case BroadcastExceptSelf:
			buffer := protocol.Encode(ctx, replyv)
			s.server.connMatrix.BroadcastExceptOne(buffer, ctx.Conn.Id())
		case BroadcastSomeone:
			buffer := protocol.Encode(ctx, replyv)
			s.server.connMatrix.BroadcastSomeone(buffer, callModel.Ids)
		}

		//rpclog.Infof("args %s", string(data))
		return nil
	}
	if isAsync {
		err := s.server.gPool.Submit(func() {
			err := f()
			if err != nil {
				rpclog.Errorf("call async err:%s", err.Error())
			}
		})
		return err
	}

	return f()
}
