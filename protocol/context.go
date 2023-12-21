package protocol

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
	"sync"
)

var (
	CtxPool = sync.Pool{
		New: func() any {
			return &Context{
				H:        &Header{},
				Metadata: map[string]string{},
				Payload:  bytebufferpool.Get(),
			}
		},
	}
)

func GetCtx() *Context {
	return CtxPool.Get().(*Context)
}
func PutCtx(ctx *Context) {
	ctx.Reset()
	CtxPool.Put(ctx)
}

type Context struct {
	H             *Header
	Payload       *bytebufferpool.ByteBuffer
	Conn          gnet.Conn
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	//SerializeType uint8
	MsgSeq uint64 //reserved field
	Ctx    context.Context
	GPool  *ants.Pool
}

func (c *Context) Reset() {
	c.H.SerializeType = 0
	c.H.MagicNumber = 0
	c.H.HeartBeat = 0
	c.H.Version = 0
	c.Payload = nil
	c.Conn = nil
	c.ServicePath = ""
	c.ServiceMethod = ""
	c.Metadata = nil
	c.MsgSeq = 0
	c.Ctx = nil
	c.GPool = nil
}
