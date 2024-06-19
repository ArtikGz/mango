package s2c

import (
	"bytes"
	dt "mango/src/network/datatypes"
)

type SystemChatMessage struct {
	Content string // The text message
	Overlay dt.Boolean

	ShouldBroadcast bool
}

func (pk SystemChatMessage) Bytes() []byte {
	contentJson := dt.String("{\"text\": \"" + pk.Content + "\"}")
	return bytes.Join([][]byte{
		dt.VarInt(0x64).Bytes(),
		contentJson.Bytes(),
		pk.Overlay.Bytes(),
	}, nil)
}

func (pk SystemChatMessage) Broadcast() bool {
	return pk.ShouldBroadcast
}
