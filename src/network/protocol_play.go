package network

import (
	"bytes"
	"io"
	"mango/src/logger"
	"mango/src/managers"
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
		// TODO: Hacer validaciones con respecto a la distancia entre el bloque y el jugador
		// TODO: Implementar formas de minado distintas para los distintos gamemodes (Que tarde en survival y que no se pueda en adventure)
		managers.GetBlockManager().RemoveBlockAt(packet.Position.X, packet.Position.Y, packet.Position.Z)
	case c2s.ACTION_STATUS_CANCELLED_DIGGING:
	case c2s.ACTION_STATUS_FINISHED_DIGGING:
	}

	logger.Debug("%+v", packet)
}
