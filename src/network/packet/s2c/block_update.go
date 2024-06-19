package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type BlockUpdate struct {
	Location dt.Position
	BlockId  dt.VarInt
}

func (pk BlockUpdate) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x0A).Bytes(),
		pk.Location.Bytes(),
		pk.BlockId.Bytes(),
	}, nil)
}

func (pk BlockUpdate) Broadcast() bool {
	// We assume that we always want to broadcast a block update
	return true
}
