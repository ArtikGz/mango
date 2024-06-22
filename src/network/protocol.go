package network

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

func HandlePacket(username string, state Protocol, packet []byte) ([]Packet, error) {
	switch state {
	case SHAKE:
		return HandleHandshakePacket(packet)
	case STATUS:
		return HandleStatusPacket(packet)
	case LOGIN:
		return HandleLoginPacket(packet)
	case PLAY:
		return HandlePlayPacket(username, packet)
	default:
		return nil, nil
	}
}
