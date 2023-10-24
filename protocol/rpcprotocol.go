package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/cat3306/gnetrpc/component"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
	"strings"
)

// 0                     4                       12          14
// +---------------------+-----------------------+-----------+
// |   payload len       |       seq             | code type |
// +---------------------+-----------------------+-----------+
// | path and method len |  path@method       			 	 |
// +---------------------------------------------------------+                                   		            			 +
// |                  payload        		             	 |
// +                                   						 |
// |                                  						 |
// +---------------------------------------------------------+

const (
	payloadLen       = uint32(4)
	seqLen           = uint32(8)
	pathMethodLen    = uint32(4)
	serializeTypeLen = uint32(2)

	maxBufferCap = 1 << 24 //16M
)

var (
	packetEndian        = binary.LittleEndian
	ErrIncompletePacket = errors.New("incomplete packet")
	ErrTooLargePacket   = errors.New("too large packet")
	ErrDiscardedPacket  = errors.New("discarded not equal msg len")
)

func Decode(c gnet.Conn) (*component.Context, error) {

	headerLen := int(payloadLen + seqLen + serializeTypeLen + pathMethodLen)
	headerBuffer, err := c.Peek(headerLen)
	if err != nil {
		return nil, err
	}
	if len(headerBuffer) < headerLen {
		return nil, ErrIncompletePacket
	}
	payloadLength := packetEndian.Uint32(headerBuffer[:payloadLen])

	seq := packetEndian.Uint64(headerBuffer[payloadLen : payloadLen+seqLen])

	serializeType := packetEndian.Uint16(headerBuffer[payloadLen+seqLen : payloadLen+seqLen+serializeTypeLen])

	pathMethodLength := int(packetEndian.Uint32(headerBuffer[payloadLen+seqLen+serializeTypeLen : payloadLen+seqLen+serializeTypeLen+pathMethodLen]))

	msgLen := headerLen + int(payloadLength) + pathMethodLength
	if msgLen > maxBufferCap {
		return nil, ErrTooLargePacket
	}
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	msgBuffer, err := c.Peek(msgLen)
	if err != nil {
		return nil, err
	}
	discarded, err := c.Discard(msgLen)
	if err != nil {
		return nil, err
	}
	if discarded != msgLen {
		rpclog.Errorf("discarded")
		return nil, ErrDiscardedPacket
	}
	methodBuffer := msgBuffer[headerLen : headerLen+pathMethodLength]
	servicePathAndMethod := strings.Split(util.BytesToString(methodBuffer), "@")
	if len(servicePathAndMethod) != 2 {

	}
	servicePath := servicePathAndMethod[0]
	method := servicePathAndMethod[1]
	buffer := bytebufferpool.Get()
	_, _ = buffer.Write(msgBuffer[headerLen+pathMethodLength:])
	ctx := &component.Context{
		ServiceMethod: method,
		ServicePath:   servicePath,
		Payload:       buffer,
		SerializeType: serializeType,
		Conn:          c,
		Seq:           seq,
	}
	return ctx, nil
}
func Encode(ctx *component.Context, v interface{}) *bytebufferpool.ByteBuffer {
	if v == nil {
		panic("v nil")
	}
	var (
		payload []byte
		err     error
	)
	if tmp, ok := v.([]byte); ok {
		payload = tmp
	} else {
		payload, err = GameCodec(SerializeType(ctx.SerializeType)).Marshal(v)
		if err != nil {
			panic(err)
		}
	}
	buffer := bytebufferpool.Get()
	headBuffer := make([]byte, int(payloadLen+seqLen+serializeTypeLen+pathMethodLen))
	packetEndian.PutUint32(headBuffer, uint32(len(payload)))
	packetEndian.PutUint64(headBuffer[payloadLen:], ctx.Seq)
	packetEndian.PutUint16(headBuffer[payloadLen+seqLen:], ctx.SerializeType)
	methodStr := ctx.ServicePath + "@" + ctx.ServiceMethod
	packetEndian.PutUint32(headBuffer[payloadLen+seqLen+serializeTypeLen:], uint32(len(methodStr)))

	_, _ = buffer.Write(headBuffer)
	_, _ = buffer.Write(util.StringToBytes(methodStr))
	_, _ = buffer.Write(payload)
	return buffer
}
