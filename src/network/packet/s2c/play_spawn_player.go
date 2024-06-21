package s2c

import (
	"bytes"
	"encoding/binary"
	dt "mango/src/network/datatypes"
)

type PlaySpawnPlayer struct {
	EntityId dt.VarInt
	UUID     []byte
	X        dt.Double
	Y        dt.Double
	Z        dt.Double
	Yaw      dt.UByte
	Pitch    dt.UByte

	ShouldBroadcast bool
}

func (pk PlaySpawnPlayer) Bytes() []byte {
	uuid1 := binary.BigEndian.Uint64(pk.UUID[:8])
	uuid2 := binary.BigEndian.Uint64(pk.UUID[8:])

	return bytes.Join([][]byte{
		dt.VarInt(0x03).Bytes(),
		pk.EntityId.Bytes(),
		dt.Long(uuid1).Bytes(), dt.Long(uuid2).Bytes(),
		pk.X.Bytes(),
		pk.Y.Bytes(),
		pk.Z.Bytes(),
		pk.Yaw.Bytes(),
		pk.Pitch.Bytes(),
	}, nil)
}

func (pk PlaySpawnPlayer) Broadcast() bool {
	return pk.ShouldBroadcast
}
