package network

import "io"

type IncomingPacket interface {
	ReadPacket(io.Reader)
}

type OutgoingPacket interface {
	Bytes() []byte
	Broadcast() bool
}

type Packet interface {
}
