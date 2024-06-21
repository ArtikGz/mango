package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type SetPlayerPositionAndRotation struct {
	X        dt.Double
	Y        dt.Double
	Z        dt.Double
	Yaw      dt.UByte
	Pitch    dt.UByte
	OnGround dt.Boolean
}

func ReadSetPlayerPositionAndRotationPacket(r io.Reader) (*SetPlayerPositionAndRotation, error) {
	var pk SetPlayerPositionAndRotation

	if _, err := pk.X.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Y.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Z.ReadFrom(r); err != nil {
		return nil, err
	}
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
