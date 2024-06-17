package network

import (
	"bytes"
	"mango/src/network/packet/c2s"
)

func HandleHandshakePacket(data *[]byte) []Packet {
	reader := bytes.NewReader(*data)

	var handshake c2s.Handshake
	handshake.ReadPacket(reader)

	return []Packet{handshake}
}
