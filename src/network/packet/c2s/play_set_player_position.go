package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type SetPlayerPosition struct {
	X        dt.Double
	Y        dt.Double
	Z        dt.Double
	OnGround dt.Boolean
}

func ReadSetPlayerPositionPacket(r io.Reader) (*SetPlayerPosition, error) {
	var pk SetPlayerPosition

	if _, err := pk.X.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Y.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Z.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.OnGround.ReadFrom(r); err != nil {
		return nil, err
	}
	return &pk, nil
}
