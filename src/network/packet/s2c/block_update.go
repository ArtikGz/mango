package s2c

import (
	dt "mango/src/network/datatypes"
	"mango/src/network/packet"
)

type BlockUpdate struct {
	Header   packet.PacketHeader
	Location dt.Position
	BlockId  dt.VarInt
}

func (pk BlockUpdate) Bytes() []byte {
	pk.Header.PacketID = 0x0A
	var data []byte

	data = append(data, pk.Location.Bytes()...)
	data = append(data, pk.BlockId.Bytes()...)

	pk.Header.WriteHeader(&data)

	return data
}

func (pk BlockUpdate) Broadcast() bool {
	// We assume that we always want to broadcast a block update
	return true
}
