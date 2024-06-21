package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type PlaySynchronizePlayerPosition struct {
	X          dt.Double
	Y          dt.Double
	Z          dt.Double
	Yaw        dt.Float
	Pitch      dt.Float
	Flags      uint8
	TeleportID dt.VarInt
}

func (pk PlaySynchronizePlayerPosition) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x3C).Bytes(),
		pk.X.Bytes(),
		pk.Y.Bytes(),
		pk.Z.Bytes(),
		pk.Yaw.Bytes(),
		pk.Pitch.Bytes(),
		{pk.Flags},
		pk.TeleportID.Bytes(),
	}, nil)
}

func (pk PlaySynchronizePlayerPosition) Broadcast() bool {
	return false
}
