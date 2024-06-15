package network

import (
	"bytes"
	"io"
	"mango/src/config"
	"mango/src/network/packet"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
	"net"
)

func HandleStatusPacket(conn *net.TCPConn, data *[]byte) []Packet {
	reader := bytes.NewReader(*data)

	var header packet.PacketHeader
	header.ReadHeader(reader)

	reader.Seek(0, io.SeekStart)

	switch header.PacketID {
	case 0x00: // status packet
		var statusRequest c2s.StatusRequest
		statusRequest.ReadPacket(reader)

		var statusResponse s2c.StatusResponse
		statusResponse.StatusData.Protocol = uint16(config.Protocol())

		return []Packet{statusResponse}

	case 0x01: // ping packet
		var ping c2s.PingRequest
		ping.ReadPacket(reader)

		var pong s2c.PingResponse
		pong.Timestamp = ping.Timestamp

		return []Packet{pong}
	}

	return nil
}
