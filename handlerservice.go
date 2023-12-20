package gnetrpc

import (
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/ants/v2"
	"reflect"
	"runtime/debug"
)

type Handler func(ctx *protocol.Context)
type AsyncHandler func(ctx *protocol.Context, tag struct{})
type HandlerSet struct {
	set      map[string]Handler
	asyncSet map[string]AsyncHandler
}

func NewHandlerSet() *HandlerSet {
	return &HandlerSet{
		set:      map[string]Handler{},
		asyncSet: map[string]AsyncHandler{},
	}
}
func (h *HandlerSet) ExecuteHandler(ctx *protocol.Context, gPool *ants.Pool) error {
	var err error
	key := util.JoinServiceMethod(ctx.ServicePath, ctx.ServiceMethod)
	handler, ok := h.set[key]
	if ok {
		handler(ctx)
		return nil
	}
	asHandler, ok := h.asyncSet[key]
	if !ok {
		return NotFoundMethod
	}
	err = gPool.Submit(func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				err := fmt.Errorf("[ExecuteHandler internal error] service: %s, method: %s,err:%v stack: %s", ctx.ServiceMethod, ctx.ServicePath, r, util.BytesToString(stack))
				rpclog.Error(err)
			}
		}()
		asHandler(ctx, struct{}{})
	})
	return err
}
func (h *HandlerSet) Register(v IService, isPrint bool, name ...string) {
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

	for i := 0; i < typ.NumMethod(); i++ {
		mName := typ.Method(i).Name
		if !util.IsExported(mName) {
			continue
		}
		f, ok := value.Method(i).Interface().(func(ctx *protocol.Context))
		if ok {
			h.set[util.JoinServiceMethod(sName, mName)] = f
			if isPrint {
				rpclog.Info(fmt.Sprintf("registered [%s.%s]", sName, mName))
			}
			continue
		}
		af, ok := value.Method(i).Interface().(func(ctx *protocol.Context, tag struct{}))
		if ok {
			h.asyncSet[util.JoinServiceMethod(sName, mName)] = af
			if isPrint {
				rpclog.Info(fmt.Sprintf("registered [%s.go_%s]", sName, mName))
			}
		}

	}
}
