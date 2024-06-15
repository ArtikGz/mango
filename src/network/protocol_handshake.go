package network

import (
	"bytes"
	"mango/src/network/packet/c2s"
	"net"
)

func HandleHandshakePacket(conn *net.TCPConn, data *[]byte) []Packet {
	reader := bytes.NewReader(*data)

	var handshake c2s.Handshake
	handshake.ReadPacket(reader)

	return []Packet{handshake}
}
