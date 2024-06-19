package network

import (
	"bytes"
	"io"
	"mango/src/logger"
	"mango/src/managers"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
)

func HandlePlayPacket(data []byte) ([]Packet, error) {
	r := bytes.NewReader(data)

	pid, _, err := dt.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	logger.Debug("PLAY packet ID: %d", pid)

	switch pid {
	// Player Action
	case 0x1d:
		return handlePlayerAction(r)
	}

	return nil, nil
}

func handlePlayerAction(r io.Reader) ([]Packet, error) {
	pk, err := c2s.ReadPlayerActionPacket(r)
	if err != nil {
		return nil, err
	}

	switch pk.Status {
	case c2s.ACTION_STATUS_STARTED_DIGGING:
		// TODO: Hacer validaciones con respecto a la distancia entre el bloque y el jugador
		// TODO: Implementar formas de minado distintas para los distintos gamemodes (Que tarde en survival y que no se pueda en adventure)
		managers.GetBlockManager().RemoveBlockAt(pk.Position.X, pk.Position.Y, pk.Position.Z)

		return []Packet{s2c.BlockUpdate{
			Location: pk.Position,
			BlockId:  0,
		}}, nil
	case c2s.ACTION_STATUS_CANCELLED_DIGGING:
	case c2s.ACTION_STATUS_FINISHED_DIGGING:
	default:
	}

	return nil, nil
}
