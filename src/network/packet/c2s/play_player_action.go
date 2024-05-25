package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type ActionStatus int32

const (
	ACTION_STATUS_STARTED_DIGGING = iota
	ACTION_STATUS_CANCELLED_DIGGING
	ACTION_STATUS_FINISHED_DIGGING
	ACTION_STATUS_DROP_ITEM_STACK
	ACTION_STATUS_DROP_ITEM
	ACTION_STATUS_SHOT_ARROW_FINISH_EATING
	ACTION_STATUS_SWAP_ITEM_IN_HAND
)

type PlayerAction struct {
	Status   dt.VarInt
	Position dt.Position
	Face     dt.Byte
	Sequence dt.VarInt
}

func (pa *PlayerAction) ReadPacket(reader io.Reader) {
	pa.Status.ReadFrom(reader)
	pa.Position.ReadFrom(reader)
	pa.Face.ReadFrom(reader)
	pa.Sequence.ReadFrom(reader)
}
