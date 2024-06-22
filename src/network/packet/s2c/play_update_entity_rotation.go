package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type PlayUpdateEntityRotation struct {
	EntityId dt.VarInt
	Yaw      dt.UByte
	Pitch    dt.UByte
	OnGround dt.Boolean
}

func (pk PlayUpdateEntityRotation) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x2D).Bytes(),
		pk.EntityId.Bytes(),
		pk.Yaw.Bytes(),
		pk.Pitch.Bytes(),
		pk.OnGround.Bytes(),
	}, nil)
}

func (pk PlayUpdateEntityRotation) Broadcast() bool {
	return true
}
