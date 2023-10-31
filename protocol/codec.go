package protocol

type SerializeType uint8

const (
	CodeNone    = SerializeType(0)
	String      = SerializeType(1)
	Json        = SerializeType(2)
	ProtoBuffer = SerializeType(3)
)

var (
	codecSet = map[SerializeType]Codec{
		Json:        &jsonCodec{},
		String:      &stringCodec{},
		ProtoBuffer: &protocBufferCodec{},
	}
)

type Codec interface {
	Unmarshal([]byte, interface{}) error   //解码
	Marshal(v interface{}) ([]byte, error) //编码
	ToString() string
}

func GetCodec(t SerializeType) Codec {
	coder := codecSet[t]
	if coder == nil {
		coder = codecSet[Json]
	}
	return coder
}
