package network

import (
	"bytes"
	"errors"
	"mango/src/config"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
)

func HandleStatusPacket(data []byte) ([]Packet, error) {
	r := bytes.NewReader(data)

	pid, _, err := dt.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	switch pid {
	case 0x00:
		var statusResponse s2c.StatusResponse
		statusResponse.StatusData.Protocol = uint16(config.Protocol())

		return []Packet{statusResponse}, nil

	case 0x1:
		ping, err := c2s.ReadPingPacket(r)
		if err != nil {
			return nil, err
		}

		var pong s2c.PingResponse
		pong.Timestamp = ping.Timestamp

		return []Packet{pong}, nil

	default:
		return nil, errors.New("invalid handshake packet PID")
	}
}
