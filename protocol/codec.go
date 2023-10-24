package protocol

type SerializeType uint8

const (
	CodeNone    = SerializeType(0)
	String      = SerializeType(1)
	Json        = SerializeType(2)
	ProtoBuffer = SerializeType(3)
)

var (
	coderSet = map[SerializeType]Codec{
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

func GameCodec(t SerializeType) Codec {
	coder := coderSet[t]
	if coder == nil {
		coder = coderSet[Json]
	}
	return coder
}
