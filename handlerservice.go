package gnetrpc

import (
	"fmt"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"reflect"
	"sync"
)

type Handler func(ctx *protocol.Context) (*CallMode, error)

type HandlerSet struct {
	set map[string]Handler
	mu  sync.RWMutex
}

func NewHandlerSet() *HandlerSet {
	return &HandlerSet{
		set: map[string]Handler{},
	}
}
func (h *HandlerSet) Register(v interface{}, isPrint bool, name ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
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
		f, ok := value.Method(i).Interface().(func(ctx *protocol.Context) (*CallMode, error))
		if !ok {
			continue
		}
		h.set[sName+"@"+mName] = f
		if isPrint {
			rpclog.Info(fmt.Sprintf("registered [%s.%s]", sName, mName))
		}
	}
}
