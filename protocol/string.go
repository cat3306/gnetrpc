package protocol

import (
	"errors"
	"fmt"
	"unsafe"
)

type stringCodec struct {
	SerializeType SerializeType
}

func (r *stringCodec) ToString() string {
	return "String"
}
func (r *stringCodec) Unmarshal(b []byte, v interface{}) error {
	if vv, ok := v.(*string); ok {
		*vv = *(*string)(unsafe.Pointer(&b))
		return nil
	}
	return errors.New("v type not string")
}
func (r *stringCodec) Marshal(v interface{}) ([]byte, error) {
	if vv, ok := v.(string); ok {
		return []byte(vv), nil
	}
	return nil, fmt.Errorf("v type not string")
}
