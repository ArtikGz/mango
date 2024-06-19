package network

import (
	"bytes"
	"errors"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
)

func HandleHandshakePacket(data []byte) ([]Packet, error) {
	r := bytes.NewReader(data)

	pid, _, err := dt.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	switch pid {
	case 0x00:
		handshake, err := c2s.ReadHandshakePacket(r)
		if err != nil {
			return nil, err
		}
		return []Packet{*handshake}, nil

	default:
		return nil, errors.New("invalid handshake packet PID")
	}
}
