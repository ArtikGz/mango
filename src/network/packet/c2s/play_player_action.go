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

func ReadPlayerActionPacket(r io.Reader) (*PlayerAction, error) {
	var pk PlayerAction

	if _, err := pk.Status.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Position.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Face.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Sequence.ReadFrom(r); err != nil {
		return nil, err
	}
	return &pk, nil
}
