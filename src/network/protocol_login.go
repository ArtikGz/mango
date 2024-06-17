package network

import (
	"bytes"
	"fmt"
	"io"
	"mango/src/config"
	"mango/src/logger"
	"mango/src/managers"
	"mango/src/network/packet"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
)

func HandleLoginPacket(data *[]byte) []Packet {
	packets := make([]Packet, 0)

	reader := bytes.NewReader(*data)

	var header packet.PacketHeader
	header.ReadHeader(reader)

	reader.Seek(0, io.SeekStart)

	switch header.PacketID {
	case 0x00: // Login Start
		var loginStart c2s.LoginStart
		loginStart.ReadPacket(reader)

		if config.IsOnline() {
			// TODO: implement cypher and return EncryptionRequest
			logger.Error("Online mod is not yet supported, please, change online to false in '%s'.", config.GetConfigPath())
		} else { // Offline mode, return LoginSuccess
			var loginSuccess s2c.LoginSuccess
			loginSuccess.Username = loginStart.Name
			if loginStart.HasUUID {
				loginSuccess.UUID = loginStart.UUID
			}

			logger.Debug("Login Success: %+v", loginSuccess)

			// send init PLAY packets (Login (Play), Default Spawn Position, etc.)
			packets = append(packets, loginSuccess)
			packets = append(packets, onSuccessfulLogin()...)
			packets = append(packets, s2c.SystemChatMessage{
				Content:         fmt.Sprintf("[+] %s joined the server.", loginSuccess.Username),
				Overlay:         false,
				ShouldBroadcast: true,
			})
		}
	}

	return packets
}

func onSuccessfulLogin() []Packet {
	packets := make([]Packet, 0)

	packets = append(packets, s2c.LoginPlay{})
	packets = append(packets, s2c.SetDefaultSpawnPosition{})

	// send 7X7 chunk square
	for _, chunkPacket := range managers.GetBlockManager().GetChunks() {
		packets = append(packets, chunkPacket)
	}

	return packets
}
