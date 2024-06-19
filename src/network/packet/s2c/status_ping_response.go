package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type PingResponse struct {
	Timestamp dt.Long
}

func (pk PingResponse) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x01).Bytes(),
		pk.Timestamp.Bytes(),
	}, nil)
}

func (pk PingResponse) Broadcast() bool {
	return false
}
