package c2s

import (
	"io"
	dt "mango/src/network/datatypes"
)

type Handshake struct {
	Protocol  dt.VarInt
	Address   dt.String
	Port      dt.UShort
	NextState dt.VarInt
}

func ReadHandshakePacket(r io.Reader) (*Handshake, error) {
	var pk Handshake
	if _, err := pk.Protocol.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Address.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.Port.ReadFrom(r); err != nil {
		return nil, err
	}
	if _, err := pk.NextState.ReadFrom(r); err != nil {
		return nil, err
	}
	return &pk, nil
}
