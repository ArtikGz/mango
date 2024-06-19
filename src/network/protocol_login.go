package network

import (
	"bytes"
	"fmt"
	"mango/src/config"
	"mango/src/logger"
	"mango/src/managers"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
)

func HandleLoginPacket(data []byte) ([]Packet, error) {
	var packets []Packet

	r := bytes.NewReader(data)

	pid, _, err := dt.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	switch pid {
	case 0x00: // Login Start

		loginStart, err := c2s.ReadLoginStartPacket(r)
		if err != nil {
			return nil, err
		}

		if config.IsOnline() {
			// TODO: implement cypher and return EncryptionRequest
			logger.Error("Online mode is not yet supported, please, change online to false in '%s'.", config.GetConfigPath())
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

	return packets, nil
}

func onSuccessfulLogin() []Packet {
	packets := []Packet{
		s2c.LoginPlay{},
		s2c.SetDefaultSpawnPosition{},
	}

	// send 7X7 chunk square
	for _, chunkPacket := range managers.GetBlockManager().GetChunks() {
		packets = append(packets, chunkPacket)
	}

	return packets
}
