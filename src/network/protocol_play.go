package network

import (
	"bytes"
	"io"
	"mango/src/logger"
	"mango/src/managers"
	"mango/src/network/packet"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
)

func HandlePlayPacket(data *[]byte) []Packet {
	reader := bytes.NewReader(*data)

	var header packet.PacketHeader
	header.ReadHeader(reader)

	logger.Debug("PLAY packet ID: %d", header.PacketID)

	switch header.PacketID {
	// Player Action
	case 0x1d:
		return handlePlayerAction(reader)
	}

	return nil
}

func handlePlayerAction(reader io.Reader) []Packet {
	var pk c2s.PlayerAction
	pk.ReadPacket(reader)

	switch pk.Status {
	case c2s.ACTION_STATUS_STARTED_DIGGING:
		// TODO: Hacer validaciones con respecto a la distancia entre el bloque y el jugador
		// TODO: Implementar formas de minado distintas para los distintos gamemodes (Que tarde en survival y que no se pueda en adventure)
		managers.GetBlockManager().RemoveBlockAt(pk.Position.X, pk.Position.Y, pk.Position.Z)

		return []Packet{s2c.BlockUpdate{
			Location: pk.Position,
			BlockId:  0,
		}}
	case c2s.ACTION_STATUS_CANCELLED_DIGGING:
	case c2s.ACTION_STATUS_FINISHED_DIGGING:
	default:
	}

	return nil
}
