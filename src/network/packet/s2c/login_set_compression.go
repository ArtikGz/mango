package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type SetCompressionPacket struct {
	Threshold dt.VarInt
}

func (pk SetCompressionPacket) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x03).Bytes(),
		pk.Threshold.Bytes(),
	}, nil)
}

func (pk SetCompressionPacket) Broadcast() bool {
	return false
}
