package protocol

import "errors"

var (
	magicNumberInvalid = errors.New("magic number invalid")
	versionInvalid     = errors.New("version number invalid")
)

const (
	MagicNumber byte = 0xFF
	Version     byte = 0x01
)

type Header struct {
	MagicNumber   byte
	Version       byte
	HeartBeat     byte //reserved field
	SerializeType byte
}

func (h *Header) Check() error {
	if h.MagicNumber != MagicNumber {
		return magicNumberInvalid
	}
	if h.Version != Version {
		return versionInvalid
	}
	return nil
}
func (h *Header) Fill(s SerializeType) {
	h.HeartBeat = 0
	h.SerializeType = byte(s)
	h.MagicNumber = MagicNumber
	h.Version = Version
}
