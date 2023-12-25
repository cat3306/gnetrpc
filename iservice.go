package gnetrpc

import (
	"github.com/cat3306/gnetrpc/protocol"
)

type IService interface {
	Init(v ...interface{}) IService
	Alias() string
}

type BaseService struct {
}

func (b *BaseService) Alias() string {
	return ""
}
func (b *BaseService) Middleware(...Handler) string {
	return ""
}
func (b *BaseService) Init(v ...interface{}) IService {
	return b
}

// Handler sync handler
func (b *BaseService) Handler(ctx *protocol.Context) {

}

// AsyncHandler Async handler
func (b *BaseService) AsyncHandler(ctx *protocol.Context, tag struct{}) {

}

// Handler0 sync handler
func (b *BaseService) Handler0(ctx *protocol.Context, req *struct{}, rsp *struct{}) *CallMode {
	return CallSelf()
}

// AsyncHandler0 Async handler
func (b *BaseService) AsyncHandler0(ctx *protocol.Context, req *struct{}, rsp *struct{}, tag struct{}) *CallMode {
	return CallSelf()
}

// PreHandler0 pre execute a smaller sort value is executed first
func (b *BaseService) PreHandler0() (f func(ctx *protocol.Context), sort int) {
	sort = 0
	f = func(ctx *protocol.Context) {
	}
	return
}

// PreHandler1 pre execute
func (b *BaseService) PreHandler1() (f func(ctx *protocol.Context), sort int) {
	sort = 1
	f = func(ctx *protocol.Context) {
	}
	return
}

// PreHandler2 pre execute
func (b *BaseService) PreHandler2() (f func(ctx *protocol.Context), sort int) {
	sort = 2
	f = func(ctx *protocol.Context) {
	}
	return
}

// PreHandler3 pre execute
func (b *BaseService) PreHandler3() (f func(ctx *protocol.Context), sort int) {
	sort = 3
	f = func(ctx *protocol.Context) {
	}
	return
}
