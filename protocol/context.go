package protocol

import (
	"context"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
)

type Context struct {
	H             *Header
	Payload       *bytebufferpool.ByteBuffer
	Conn          gnet.Conn
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	//SerializeType uint8
	MsgSeq uint64
	Ctx    context.Context
}
