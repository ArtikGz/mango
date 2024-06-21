package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type SetPlayerRotation struct {
	Yaw      dt.Float
	Pitch    dt.Float
	OnGround dt.Boolean
}

func ReadSetPlayerRotationPacket(r io.Reader) (*SetPlayerRotation, error) {
	var pk SetPlayerRotation

	if _, err := pk.Yaw.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Pitch.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.OnGround.ReadFrom(r); err != nil {
		return nil, err
	}
	return &pk, nil
}
