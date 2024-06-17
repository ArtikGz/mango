package s2c

import (
	dt "mango/src/network/datatypes"
	"mango/src/network/packet"
)

type SystemChatMessage struct {
	Header  packet.PacketHeader
	Content string // The text message
	Overlay dt.Boolean

	ShouldBroadcast bool
}

func (pk SystemChatMessage) Bytes() []byte {
	pk.Header.PacketID = 0x64
	var data []byte

	contentJson := dt.String("{\"text\": \"" + pk.Content + "\"}")

	data = append(data, contentJson.Bytes()...)
	data = append(data, pk.Overlay.Bytes()...)

	pk.Header.WriteHeader(&data)

	return data
}

func (pk SystemChatMessage) Broadcast() bool {
	return pk.ShouldBroadcast
}
