package s2c

import (
	"bytes"
	"encoding/binary"
	dt "mango/src/network/datatypes"
	"math/rand"
)

type LoginSuccess struct {
	Username dt.String
	UUID     []byte
}

func (pk LoginSuccess) Bytes() []byte {
	var uuid []byte
	if pk.UUID != nil && len(pk.UUID) == 16 {
		uuid = pk.UUID
	} else {
		uuid = make([]byte, 16)
		binary.BigEndian.PutUint64(uuid[:8], rand.Uint64())
		binary.BigEndian.PutUint64(uuid[8:], rand.Uint64())
	}

	return bytes.Join([][]byte{
		dt.VarInt(0x02).Bytes(),
		uuid,
		pk.Username.Bytes(),
		{0x00},
	}, nil)
}

func (pk LoginSuccess) Broadcast() bool {
	return false
}
