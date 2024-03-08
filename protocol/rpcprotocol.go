package protocol

import (
	"encoding/binary"
	"errors"
	"strings"

	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
)

// header(4 byte)+msgSeq(8 byte)+pathMethodLen(4 byte)+metaDataLen(4 byte)+payloadLen(4 byte)
// path@method(n byte)+metaData(n byte)+payload(n byte)
const (
	headerLen     = uint32(4)
	msgSeqLen     = uint32(8)
	pathMethodLen = uint32(4)
	metaDataLen   = uint32(4)
	payloadLen    = uint32(4)
	maxBufferCap  = uint32(1 << 24) //16M
)

var (
	packetEndian         = binary.LittleEndian
	ErrIncompletePacket  = errors.New("incomplete packet")
	ErrTooLargePacket    = errors.New("too large packet")
	ErrDiscardedPacket   = errors.New("discarded not equal msg len")
	ErrInvalidMethodPath = errors.New("invalid method path")
)

func Decode(c gnet.Conn) (*Context, error) {
	fixedLen := headerLen + msgSeqLen + pathMethodLen + metaDataLen + payloadLen
	fixedBuffer, _ := c.Peek(int(fixedLen))
	if len(fixedBuffer) < int(fixedLen) {
		return nil, ErrIncompletePacket
	}
	header := fixedBuffer[:headerLen]
	msgSeq := packetEndian.Uint64(fixedBuffer[headerLen : headerLen+msgSeqLen])
	pathMethodLength := packetEndian.Uint32(fixedBuffer[headerLen+msgSeqLen : headerLen+msgSeqLen+pathMethodLen])
	metaDataLength := packetEndian.Uint32(fixedBuffer[headerLen+msgSeqLen+pathMethodLen : headerLen+msgSeqLen+pathMethodLen+metaDataLen])
	payloadLength := packetEndian.Uint32(fixedBuffer[headerLen+msgSeqLen+pathMethodLen+metaDataLen : headerLen+msgSeqLen+pathMethodLen+metaDataLen+payloadLen])

	packetLen := fixedLen + pathMethodLength + metaDataLength + payloadLength
	if packetLen > maxBufferCap {
		return nil, ErrTooLargePacket
	}

	if c.InboundBuffered() < int(packetLen) {
		return nil, ErrIncompletePacket
	}

	packetBuffer, err := c.Peek(int(packetLen))
	if err != nil {
		return nil, err
	}
	discarded, err := c.Discard(int(packetLen))
	if err != nil {
		return nil, err
	}
	if discarded != int(packetLen) {
		rpclog.Errorf("discarded")
		return nil, ErrDiscardedPacket
	}
	methodData := packetBuffer[fixedLen : fixedLen+pathMethodLength]
	metaData := packetBuffer[fixedLen+pathMethodLength : fixedLen+pathMethodLength+metaDataLength]
	payload := packetBuffer[fixedLen+pathMethodLength+metaDataLength:]

	pathAndMethod := strings.Split(string(methodData), "@") //fix util.BytesToString()
	if len(pathAndMethod) != 2 {
		return nil, ErrInvalidMethodPath
	}
	path := pathAndMethod[0]
	method := pathAndMethod[1]
	buffer := bytebufferpool.Get()
	_, _ = buffer.Write(payload)
	ctx := GetCtx()
	//h := Header{
	//	MagicNumber:   header[0],
	//	Version:       header[1],
	//	HeartBeat:     header[2],
	//	SerializeType: header[3],
	//}
	ctx.H.MagicNumber = header[0]
	ctx.H.Version = header[1]
	ctx.H.HeartBeat = header[2]
	ctx.H.SerializeType = header[3]
	ctx.ServiceMethod = method
	ctx.ServicePath = path
	ctx.Payload = buffer
	ctx.Conn = c
	ctx.MsgSeq = msgSeq
	if len(metaData) != 0 {
		codec := GetCodec(Json)
		err = codec.Unmarshal(metaData, &ctx.Metadata)
		if err != nil {
			return nil, err
		}
	}
	err = ctx.H.Check()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}
func Encode(ctx *Context, v interface{}) *bytebufferpool.ByteBuffer {
	if ctx.H == nil {
		panic("encode header nil")
	}
	var (
		payload  []byte
		err      error
		metaData []byte
	)
	codec := GetCodec(SerializeType(ctx.H.SerializeType))
	if ctx.Metadata != nil && len(ctx.Metadata) != 0 {
		pc := GetCodec(Json)
		metaData, err = pc.Marshal(ctx.Metadata)
	}
	if err != nil {
		panic(err)
	}
	if v != nil {
		if tmp, ok := v.([]byte); ok {
			payload = tmp
		} else {
			payload, err = codec.Marshal(v)
			if err != nil {
				panic(err)
			}
		}
	}
	buffer := bytebufferpool.Get()
	fixedLen := headerLen + msgSeqLen + pathMethodLen + metaDataLen + payloadLen
	fixedBuffer := make([]byte, int(fixedLen))
	fixedBuffer[0] = ctx.H.MagicNumber
	fixedBuffer[1] = ctx.H.Version
	fixedBuffer[2] = ctx.H.HeartBeat
	fixedBuffer[3] = ctx.H.SerializeType
	//copy(fixedBuffer[:headerLen], ctx.H[:])
	packetEndian.PutUint64(fixedBuffer[headerLen:], ctx.MsgSeq)
	//packetEndian.PutUint16(fixedBuffer[payloadLen+seqLen:], ctx.SerializeType)
	methodStr := util.JoinServiceMethod(ctx.ServicePath, ctx.ServiceMethod)
	packetEndian.PutUint32(fixedBuffer[headerLen+msgSeqLen:], uint32(len(methodStr)))
	packetEndian.PutUint32(fixedBuffer[headerLen+msgSeqLen+pathMethodLen:], uint32(len(metaData)))
	packetEndian.PutUint32(fixedBuffer[headerLen+msgSeqLen+pathMethodLen+metaDataLen:], uint32(len(payload)))

	_, _ = buffer.Write(fixedBuffer)
	_, _ = buffer.Write(util.StringToBytes(methodStr))
	_, _ = buffer.Write(metaData)
	_, _ = buffer.Write(payload)
	return buffer
}
