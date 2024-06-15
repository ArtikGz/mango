package network

import "io"

type IncommingPacket interface {
	ReadPacket(io.Reader)
}

type OutgoingPacket interface {
	Bytes() []byte
}

type Packet interface {
}
