package gnetrpc

import (
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/ants/v2"
	"reflect"
	"runtime/debug"
	"sort"
)

type Handler func(ctx *protocol.Context)
type AsyncHandler func(ctx *protocol.Context, tag struct{})
type PreHandler func() (func(ctx *protocol.Context), int)

type PreFunc struct {
	f     Handler
	order int
}
type PreFuncList []PreFunc

func (p PreFuncList) Less(i, j int) bool {
	return p[i].order < p[j].order
}
func (p PreFuncList) Len() int {
	return len(p)
}
func (p PreFuncList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PreFuncList) Call(ctx *protocol.Context) {
	for _, pre := range p {
		pre.f(ctx)
	}
}

type HandlerSet struct {
	set        map[string]Handler
	asyncSet   map[string]AsyncHandler
	preHandler PreFuncList
}

func NewHandlerSet() *HandlerSet {
	return &HandlerSet{
		set:        map[string]Handler{},
		asyncSet:   map[string]AsyncHandler{},
		preHandler: make([]PreFunc, 0),
	}
}
func (h *HandlerSet) Call(ctx *protocol.Context, gPool *ants.Pool) error {
	var err error
	key := ctx.ServiceMethod
	handler, ok := h.set[key]
	if ok {
		// sync
		h.preHandler.Call(ctx)
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
				err := fmt.Errorf("[Call internal error] service: %s, method: %s,err:%v stack: %s", ctx.ServiceMethod, ctx.ServicePath, r, util.BytesToString(stack))
				rpclog.Error(err)
			}
		}()
		h.preHandler.Call(ctx)
		asHandler(ctx, struct{}{})
	})
	return err
}
func (h *HandlerSet) Register(v IService, isPrint bool) *HandlerSet {
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

	for i := 0; i < typ.NumMethod(); i++ {
		mName := typ.Method(i).Name
		if !util.IsExported(mName) {
			continue
		}
		f, ok := value.Method(i).Interface().(func(ctx *protocol.Context))
		if ok {
			h.set[mName] = f
			if isPrint {
				rpclog.Info(fmt.Sprintf("registered [%s.%s]", sName, mName))
			}
			continue
		}
		af, ok := value.Method(i).Interface().(func(ctx *protocol.Context, tag struct{}))
		if ok {
			h.asyncSet[mName] = af
			if isPrint {
				rpclog.Info(fmt.Sprintf("registered [%s.go@%s]", sName, mName))
			}
		}
		preFunc, ok := value.Method(i).Interface().(func() (func(ctx *protocol.Context), int))
		if ok {
			preF, order := preFunc()
			h.preHandler = append(h.preHandler, PreFunc{
				f:     preF,
				order: order,
			})
			if isPrint {
				rpclog.Info(fmt.Sprintf("registered [%s.Pre@%s]", sName, mName))
			}
		}
	}
	sort.Sort(h.preHandler)
	return h
}
