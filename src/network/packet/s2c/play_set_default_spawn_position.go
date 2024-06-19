package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type SetDefaultSpawnPosition struct {
	Location dt.Position
	Angle    dt.Float
}

func (pk SetDefaultSpawnPosition) Bytes() []byte {
	return bytes.Join([][]byte{
		dt.VarInt(0x50).Bytes(),
		pk.Location.Bytes(),
		pk.Angle.Bytes(),
	}, nil)
}

func (pk SetDefaultSpawnPosition) Broadcast() bool {
	return false
}
