package protocol

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

type protocBufferCodec struct {
	SerializeType SerializeType
}

func (p *protocBufferCodec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("type err:%T", v)
	}
	return proto.Marshal(vv)
}
func (p *protocBufferCodec) Unmarshal(bin []byte, v interface{}) error { //解码
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("type err:%T", v)
	}
	return proto.Unmarshal(bin, vv)
}
func (p *protocBufferCodec) ToString() string {
	return "ProtoBuffer"
}
