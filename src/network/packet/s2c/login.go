package s2c

import (
	dt "mango/src/network/datatypes"
	"mango/src/network/packet"
)

type LoginSuccess struct {
	Header   packet.PacketHeader
	Username dt.String
}

func (pk *LoginSuccess) Bytes() []byte {
	pk.Header.PacketID = 2
	uuid1 := dt.Long(0xEDD)
	uuid2 := dt.Long(0x1337)
	var dataBuf []byte
	dataBuf = append(dataBuf, uuid1.Bytes()...)
	dataBuf = append(dataBuf, uuid2.Bytes()...)
	dataBuf = append(dataBuf, pk.Username.Bytes()...)
	pk.Header.WriteHeader(&dataBuf)

	return dataBuf
}
