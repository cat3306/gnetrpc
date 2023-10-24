package component

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
)

type Context struct {
	Payload       *bytebufferpool.ByteBuffer
	Conn          gnet.Conn
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	SerializeType uint16
	Seq           uint64
}
