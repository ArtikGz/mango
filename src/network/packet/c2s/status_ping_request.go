package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type PingRequest struct {
	Timestamp dt.Long
}

func ReadPingPacket(r io.Reader) (*PingRequest, error) {
	var pk PingRequest
	if _, err := pk.Timestamp.ReadFrom(r); err != nil {
		return nil, err
	}

	return &pk, nil
}
