package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type PlayUpdateEntityPosition struct {
	EntityId dt.VarInt
	DeltaX   dt.Short
	DeltaY   dt.Short
	DeltaZ   dt.Short
	OnGround dt.Boolean
}

func (pk PlayUpdateEntityPosition) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x2B).Bytes(),
		pk.EntityId.Bytes(),
		pk.DeltaX.Bytes(),
		pk.DeltaY.Bytes(),
		pk.DeltaZ.Bytes(),
		pk.OnGround.Bytes(),
	}, nil)
}

func (pk PlayUpdateEntityPosition) Broadcast() bool {
	return true
}
