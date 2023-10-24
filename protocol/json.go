package protocol

import "encoding/json"

type jsonCodec struct {
	SerializeType SerializeType
}

func (j *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (j *jsonCodec) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
func (j *jsonCodec) ToString() string {
	return "Json"
}
