package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type LoginSuccess struct {
	Username dt.String
	UUID     []byte
}

func (pk LoginSuccess) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x02).Bytes(),
		pk.UUID,
		pk.Username.Bytes(),
		{0x00},
	}, nil)
}

func (pk LoginSuccess) Broadcast() bool {
	return false
}
