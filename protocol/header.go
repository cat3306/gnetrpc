package protocol

type Header struct {
	MagicNumber   byte
	Version       byte
	HeartBeat     byte
	SerializeType byte
}
