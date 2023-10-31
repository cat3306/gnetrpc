package gnetrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"reflect"
	"sync"
)

type ServiceSet struct {
	set map[string]*service
	mu  sync.RWMutex
}

func NewServiceSet() *ServiceSet {
	return &ServiceSet{
		set: map[string]*service{},
	}
}

func (s *ServiceSet) suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver, context.Context, *args, *reply.
		if mtype.NumIn() != 4 {
			//if reportErr {
			//	rpclog.Debug("method ", mname, " has wrong number of ins:", mtype.NumIn())
			//}
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if !ctxType.Implements(typeOfContext) {
			//if reportErr {
			//	rpclog.Debug("method ", mname, " must use context.Context as the first parameter")
			//}
			continue
		}

		// Second arg need not be a pointer.
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			//if reportErr {
			//	rpclog.Info(mname, " parameter type not exported: ", argType)
			//}
			continue
		}
		// Third arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			//if reportErr {
			//	rpclog.Info("method", mname, " reply type not a pointer:", replyType)
			//}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			//if reportErr {
			//	rpclog.Info("method", mname, " reply type not exported:", replyType)
			//}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			//if reportErr {
			//	rpclog.Info("method", mname, " has wrong number of outs:", mtype.NumOut())
			//}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			//if reportErr {
			//	rpclog.Info("method", mname, " returns ", returnType.String(), " not error")
			//}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}

		// init pool for reflect.Type of args and reply
		reflectTypePools.Init(argType)
		reflectTypePools.Init(replyType)
	}
	return methods
}
func (s *ServiceSet) Register(v interface{}, isPrint bool, name ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	tmpService := &service{
		name:     sName,
		value:    value,
		typ:      typ,
		method:   s.suitableMethods(typ, true),
		function: nil,
	}
	if isPrint {
		for _, m := range tmpService.method {
			rpclog.Info(fmt.Sprintf("registered [%s.%s]", tmpService.name, m.method.Name))
		}
	}
	s.set[tmpService.name] = tmpService
}

func (s *ServiceSet) Call(ctx *protocol.Context) {
	servicePath := ctx.ServicePath
	method := ctx.ServiceMethod
	tmpService := s.set[servicePath]
	if tmpService == nil {
		err := errors.New("gnetrpc: can't find service " + servicePath)
		rpclog.Error(err.Error())
		return
	}
	mtype := tmpService.method[method]
	if mtype == nil {
		err := errors.New("rpcx: can't find method " + method)
		rpclog.Error(err.Error())
		return
	}
	replyv := reflectTypePools.Get(mtype.ReplyType)
	argv := reflectTypePools.Get(mtype.ArgType)
	codec := protocol.GetCodec(protocol.SerializeType(ctx.SerializeType))
	if codec == nil {
		err := errors.New("rpcx: can't find method " + method)
		rpclog.Error(err.Error())
		return
	}
	err := codec.Unmarshal(ctx.Payload.Bytes(), argv)
	if err != nil {
		err := errors.New("rpcx: can't find method " + method)
		rpclog.Error(err.Error())
		return
	}
	err = tmpService.call(context.Background(), mtype, reflect.ValueOf(argv), reflect.ValueOf(replyv))
	return
}
