package tcp

import (
	"mango/src/network"
	"net"
)

type ( // events
	ClientConnectEvent struct {
		conn *net.TCPConn
	}

	ClientDisconnectEvent struct {
		client *TcpClient
	}

	ClientCrashEvent struct {
		client *TcpClient
		err    error
	}
)

type (
	BroadcastPacketEvent struct {
		packet network.OutgoingPacket
	}
)
