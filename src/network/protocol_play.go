package network

import (
	"bytes"
	"io"
	"mango/src/logger"
	"mango/src/managers"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
	"mango/src/world"
)

func HandlePlayPacket(ctx PacketContext, data []byte) ([]Packet, error) {
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
	case 0x14:
		return handleSetPlayerPosition(ctx, r)
	case 0x15:
		return handleSetPlayerPositionAndRotation(ctx, r)
	case 0x16:
		return handleSetPlayerRotation(ctx, r)
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
		world.WorldInstance.RemoveBlock(world.BlockPos{
			X: pk.Position.X,
			Y: pk.Position.Y,
			Z: pk.Position.Z,
		})

		return []Packet{s2c.BlockUpdate{
			Location: pk.Position,
			BlockId:  0, // Air
		}}, nil
	case c2s.ACTION_STATUS_CANCELLED_DIGGING:
	case c2s.ACTION_STATUS_FINISHED_DIGGING:
	default:
	}

	return nil, nil
}

func handleSetPlayerPosition(ctx PacketContext, r io.Reader) ([]Packet, error) {
	pk, err := c2s.ReadSetPlayerPositionPacket(r)
	if err != nil {
		return nil, err
	}

	user := managers.GetUserManager().GetUser(ctx.Username())
	deltaX := calcPositionDiff(user.Position.X, float64(pk.X))
	deltaY := calcPositionDiff(user.Position.Y, float64(pk.Y))
	deltaZ := calcPositionDiff(user.Position.Z, float64(pk.Z))
	//logger.Debug("Player (EntityID: %d) position change requested: {%f, %f, %f} -> {%f, %f, %f} (Delta: {%d, %d, %d})", user.EntityId, user.Position.X, user.Position.Y, user.Position.Z, pk.X, pk.Y, pk.Z, deltaX, deltaY, deltaZ)

	user.Position.X = float64(pk.X)
	user.Position.Y = float64(pk.Y)
	user.Position.Z = float64(pk.Z)
	managers.GetUserManager().UpdateUser(*user)

	return []Packet{
		s2c.PlayUpdateEntityPosition{
			EntityId: dt.VarInt(user.EntityId),
			DeltaX:   dt.Short(deltaX),
			DeltaY:   dt.Short(deltaY),
			DeltaZ:   dt.Short(deltaZ),
			OnGround: pk.OnGround,
		},
	}, nil
}

func handleSetPlayerPositionAndRotation(ctx PacketContext, r io.Reader) ([]Packet, error) {
	pk, err := c2s.ReadSetPlayerPositionAndRotationPacket(r)
	if err != nil {
		return nil, err
	}

	user := managers.GetUserManager().GetUser(ctx.Username())
	deltaX := calcPositionDiff(user.Position.X, float64(pk.X))
	deltaY := calcPositionDiff(user.Position.Y, float64(pk.Y))
	deltaZ := calcPositionDiff(user.Position.Z, float64(pk.Z))
	// TODO: log this
	// logger.Debug("Player (EntityID: %d) position and rotation change requested: {%f, %f, %f, %f, %f} -> {%f, %f, %f, %f, %f} (Delta: {%d, %d, %d})", user.EntityId, user.Position.X, user.Position.Y, user.Position.Z, pk.X, pk.Y, pk.Z, deltaX, deltaY, deltaZ)

	user.Position.X = float64(pk.X)
	user.Position.Y = float64(pk.Y)
	user.Position.Z = float64(pk.Z)
	user.Position.Yaw = uint8(pk.Yaw)
	user.Position.Pitch = uint8(pk.Pitch)
	managers.GetUserManager().UpdateUser(*user)

	return []Packet{
		s2c.PlayUpdateEntityPositionAndRotation{
			EntityId: dt.VarInt(user.EntityId),
			DeltaX:   dt.Short(deltaX),
			DeltaY:   dt.Short(deltaY),
			DeltaZ:   dt.Short(deltaZ),
			Yaw:      pk.Yaw,
			Pitch:    pk.Pitch,
			OnGround: pk.OnGround,
		},
	}, nil
}

func handleSetPlayerRotation(ctx PacketContext, r io.Reader) ([]Packet, error) {
	pk, err := c2s.ReadSetPlayerRotationPacket(r)
	if err != nil {
		return nil, err
	}

	user := managers.GetUserManager().GetUser(ctx.Username())

	return []Packet{
		s2c.PlayUpdateEntityRotation{
			EntityId: dt.VarInt(user.EntityId),
			Yaw:      dt.UByte(pk.Yaw),
			Pitch:    dt.UByte(pk.Pitch),
			OnGround: pk.OnGround,
		},
	}, nil
}

func calcPositionDiff(prevX, newX float64) int16 {
	return int16((newX*32 - prevX*32) * 128)
}
