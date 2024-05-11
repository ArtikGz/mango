package network

import (
	"bytes"
	"io"
	"mango/src/logger"
	"mango/src/network/packet"
	"mango/src/network/packet/c2s"
)

func HandlePlayPacket(conn *Connection, data *[]byte) {
	reader := bytes.NewReader(*data)

	var header packet.PacketHeader
	header.ReadHeader(reader)

	logger.Debug("PLAY packet ID: %d", header.PacketID)

	switch header.PacketID {
	// Player Action
	case 0x1d:
		handlePlayerAction(reader)
	}
}

func handlePlayerAction(reader io.Reader) {
	var packet c2s.PlayerAction
	packet.ReadPacket(reader)

	switch packet.Status {
	case c2s.ACTION_STATUS_STARTED_DIGGING:

	}
	logger.Debug("%+v", packet)
}
