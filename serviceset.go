package gnetrpc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/ants/v2"
)

type ServiceSet struct {
	set        map[string]*Service
	connMatrix *ConnMatrix
	gPool      *ants.Pool
}

func NewServiceSet(pool *ants.Pool, connMatrix *ConnMatrix) *ServiceSet {
	return &ServiceSet{
		set:        map[string]*Service{},
		gPool:      pool,
		connMatrix: connMatrix,
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
		if mtype.NumOut() != 1 {
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfCallMode {
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
func (s *ServiceSet) Register(v IService, isPrint bool) {
	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	sName := reflect.Indirect(value).Type().Name()
	if v.Alias() != "" {
		sName = v.Alias()
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
		handlerSet:  NewHandlerSet().Register(v, isPrint),
	}
	if isPrint {
		for _, m := range tmpService.method {
			rpclog.Info(fmt.Sprintf("registered [%s.%s]", tmpService.name, m.method.Name))
		}
		for _, m := range tmpService.asyncMethod {
			rpclog.Info(fmt.Sprintf("registered [%s.go@%s]", tmpService.name, m.method.Name))
		}

	}
	s.set[tmpService.name] = tmpService

}

// 同步
func (s *ServiceSet) SyncCall(ctx *protocol.Context) error {
	servicePath := ctx.ServicePath
	method := ctx.ServiceMethod
	tmpService := s.set[servicePath]
	if tmpService == nil {
		return NotFoundService
	}
	err := tmpService.handlerSet.SyncCall(ctx)
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, NotFoundMethod) {
		return err
	}
	mType := tmpService.method[method]
	if mType == nil {
		return NotFoundMethod
	}
	return s.Call(ctx, tmpService, mType, false)
}

// 异步
func (s *ServiceSet) AsyncCall(ctx *protocol.Context) error {
	servicePath := ctx.ServicePath
	method := ctx.ServiceMethod
	tmpService := s.set[servicePath]
	if tmpService == nil {
		return NotFoundService
	}
	err := tmpService.handlerSet.AsyncCall(ctx)
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, NotFoundMethod) {
		return err
	}
	mType := tmpService.asyncMethod[method]
	if mType == nil {
		return NotFoundMethod
	}
	return s.Call(ctx, tmpService, mType, true)
}
func (s *ServiceSet) Call(ctx *protocol.Context, tmpService *Service, mType *methodType, isAsync bool) error {
	defer func() {
		if r := recover(); r != nil {
			msg := debug.Stack()
			err := fmt.Errorf("[server call internal error] service: %s, method: %s, stack: %s,err:%s", ctx.ServicePath, ctx.ServiceMethod, util.BytesToString(msg), r)
			rpclog.Error(err)
		}
	}()
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
			protocol.PutCtx(ctx)
		}()
		//maybe ctx.Payload nil
		if ctx.Payload.Len() != 0 {
			err := codec.Unmarshal(ctx.Payload.Bytes(), argv)
			if err != nil {

				return fmt.Errorf("codec.Unmarshal err:%s", err.Error())
			}
		}
		tmpService.handlerSet.preHandler.Call(ctx)
		callModel := tmpService.call(ctx, mType, reflect.ValueOf(argv), reflect.ValueOf(replyv), isAsync)
		if callModel == nil {
			return nil
		}
		switch callModel.Call {
		case None:
			return nil
		case Self:
			buffer := protocol.Encode(ctx, replyv)
			s.connMatrix.SendToConn(buffer, ctx.Conn)
			return nil
		case Broadcast:
			buffer := protocol.Encode(ctx, replyv)
			s.connMatrix.Broadcast(buffer)
		case BroadcastExceptSelf:
			buffer := protocol.Encode(ctx, replyv)
			s.connMatrix.BroadcastExceptOne(buffer, util.GetConnId(ctx.Conn))
		case BroadcastSomeone:
			buffer := protocol.Encode(ctx, replyv)
			s.connMatrix.BroadcastSomeone(buffer, callModel.Ids)
		}

		//rpclog.Infof("args %s", string(data))
		return nil
	}
	return f()
}
