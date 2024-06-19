package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type KeepAlive struct {
	KeepAliveID dt.Long
}

func (pk KeepAlive) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x23).Bytes(),
		pk.KeepAliveID.Bytes(),
	}, nil)
}

func (pk KeepAlive) Broadcast() bool {
	return false
}
