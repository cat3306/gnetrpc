package protocol

import (
	"fmt"
	"github.com/cat3306/gnetrpc/util"
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
	return fmt.Errorf("v type not string,type is %T", v)
}
func (r *stringCodec) Marshal(v interface{}) ([]byte, error) {
	switch v.(type) {
	case string:
		return util.StringToBytes(v.(string)), nil
	case *string:
		return util.StringToBytes(*v.(*string)), nil
	default:
		return nil, fmt.Errorf("v type not string,type is %T", v)
	}
}
