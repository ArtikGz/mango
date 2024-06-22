package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type PlayUpdateEntityPositionAndRotation struct {
	EntityId dt.VarInt
	DeltaX   dt.Short
	DeltaY   dt.Short
	DeltaZ   dt.Short
	Yaw      dt.UByte
	Pitch    dt.UByte
	OnGround dt.Boolean
}

func (pk PlayUpdateEntityPositionAndRotation) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x2C).Bytes(),
		pk.EntityId.Bytes(),
		pk.DeltaX.Bytes(),
		pk.DeltaY.Bytes(),
		pk.DeltaZ.Bytes(),
		pk.Yaw.Bytes(),
		pk.Pitch.Bytes(),
		pk.OnGround.Bytes(),
	}, nil)
}

func (pk PlayUpdateEntityPositionAndRotation) Broadcast() bool {
	return true
}
