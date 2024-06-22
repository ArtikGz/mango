package s2c

import (
	"bytes"
	"encoding/binary"
	dt "mango/src/network/datatypes"
)

type Action interface {
	GetActionBitMask() uint8
	Bytes() []byte
}

type ActionAddPlayer struct {
	Name       dt.String
	Properties []dt.Property
}

func (a ActionAddPlayer) Bytes() []byte {
	props := make([]byte, 0)
	for _, prop := range a.Properties {
		props = append(props, prop.Bytes()...)
	}

	return bytes.Join([][]byte{
		a.Name.Bytes(),
		dt.VarInt(len(a.Properties)).Bytes(),
		props,
	}, nil)
}

func (a ActionAddPlayer) GetActionBitMask() uint8 {
	return 0b0000_0001
}

type ActionUpdateListed struct {
	Listed dt.Boolean
}

func (a ActionUpdateListed) Bytes() []byte {
	return bytes.Join([][]byte{
		a.Listed.Bytes(),
	}, nil)
}

func (a ActionUpdateListed) GetActionBitMask() uint8 {
	return 0b0000_1000
}

type PlayerInfoUpdate struct {
	UUID            []byte
	ActionList      []Action
	ShouldBroadcast bool
}

func (pk PlayerInfoUpdate) Bytes() []byte {
	bitMask := uint8(0)
	actionList := make([]byte, 0)
	uuid1 := binary.BigEndian.Uint64(pk.UUID[:8])
	uuid2 := binary.BigEndian.Uint64(pk.UUID[8:])
	for _, a := range pk.ActionList {
		bitMask = bitMask | a.GetActionBitMask()

		actionList = append(actionList, dt.Long(uuid1).Bytes()...)
		actionList = append(actionList, dt.Long(uuid2).Bytes()...)
		actionList = append(actionList, a.Bytes()...)
	}

	// FIXME when pk.ActionList has more than one element doesn't convert properly to []byte
	return bytes.Join([][]byte{
		dt.VarInt(0x3A).Bytes(),
		{bitMask},
		dt.VarInt(len(pk.ActionList)).Bytes(),
		actionList,
	}, nil)
}

func (pk PlayerInfoUpdate) Broadcast() bool {
	return pk.ShouldBroadcast
}
