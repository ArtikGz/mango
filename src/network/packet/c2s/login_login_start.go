package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type LoginStart struct {
	Name    dt.String
	HasUUID dt.Boolean
	UUID    []byte
}

func ReadLoginStartPacket(r io.Reader) (*LoginStart, error) {
	var pk LoginStart

	if _, err := pk.Name.ReadFrom(r); err != nil {
		return nil, err
	}

	if _, err := pk.HasUUID.ReadFrom(r); err != nil {
		return nil, err
	}
	if pk.HasUUID {
		pk.UUID = make([]byte, 16)

		if _, err := r.Read(pk.UUID); err != nil {
			return nil, err
		}
	}
	return &pk, nil
}
