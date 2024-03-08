package protocol

import (
	"google.golang.org/protobuf/proto"
)

type protocBufferCodec struct {
	SerializeType SerializeType
}

func (p *protocBufferCodec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		//return nil, fmt.Errorf("type err:%T", v)
	}
	return proto.Marshal(vv)
}
func (p *protocBufferCodec) Unmarshal(bin []byte, v interface{}) error { //解码
	vv, ok := v.(proto.Message)
	if !ok {
		//return fmt.Errorf("type err:%T", v)
	}
	return proto.Unmarshal(bin, vv)
}
func (p *protocBufferCodec) ToString() string {
	return "ProtoBuffer"
}

// func (c protocBufferCodec) Marshal(i interface{}) ([]byte, error) {
// 	if m, ok := i.(proto.Marshaler); ok {
// 		return m.Marshal()
// 	}

// 	if m, ok := i.(proto.Message); ok {
// 		return proto.Marshal(m)
// 	}

// 	return nil, fmt.Errorf("%T is not a proto.Marshaler or pb.Message", i)
// }

// // Decode decodes an object from slice of bytes.
// func (c PBCodec) Decode(data []byte, i interface{}) error {
// 	if m, ok := i.(proto.Unmarshaler); ok {
// 		return m.Unmarshal(data)
// 	}

// 	if m, ok := i.(pb.Message); ok {
// 		return pb.Unmarshal(data, m)
// 	}

// 	return fmt.Errorf("%T is not a proto.Unmarshaler  or pb.Message", i)
// }
