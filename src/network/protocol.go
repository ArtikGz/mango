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

type PacketContext interface {
	Username() string
	Protocol() Protocol
}

func HandlePacket(ctx PacketContext, packet []byte) ([]Packet, error) {
	switch ctx.Protocol() {
	case SHAKE:
		return HandleHandshakePacket(packet)
	case STATUS:
		return HandleStatusPacket(packet)
	case LOGIN:
		return HandleLoginPacket(packet)
	case PLAY:
		return HandlePlayPacket(ctx, packet)
	default:
		return nil, nil
	}
}
