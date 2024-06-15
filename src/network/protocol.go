package network

import "net"

type Protocol int

const (
	SHAKE Protocol = iota
	STATUS
	LOGIN
	PLAY
)

func (p Protocol) String() string {
	switch p {
	case SHAKE:
		return "HANDSHAKE"
	case STATUS:
		return "STATUS"
	case LOGIN:
		return "LOGIN"
	case PLAY:
		return "PLAY"
	default:
		return "UNKNOWN"
	}
}

func HandlePacket(state Protocol, conn *net.TCPConn, packet *[]byte) []Packet {
	switch state {
	case SHAKE:
		return HandleHandshakePacket(conn, packet)
	case STATUS:
		return HandleStatusPacket(conn, packet)
	case LOGIN:
		return HandleLoginPacket(conn, packet)
	case PLAY:
		return HandlePlayPacket(conn, packet)
	}

	return nil
}
