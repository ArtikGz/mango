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

			entityId := managers.GetEntityManager().GenerateID()
			userPosition := managers.UserPosition{
				X:     8.5,
				Y:     -46,
				Z:     8.5,
				Yaw:   0,
				Pitch: 0,
			}
			managers.GetUserManager().AddUser(managers.User{
				EntityId: entityId,
				UUID:     loginSuccess.UUID,
				Name:     string(loginSuccess.Username),
				Position: userPosition,
			})

			// send init PLAY packets (Login (Play), Default Spawn Position, etc.)
			packets = append(packets, loginSuccess)
			packets = append(packets, onSuccessfulLogin()...)

			packets = append(packets, s2c.PlaySynchronizePlayerPosition{
				X:     dt.Double(userPosition.X),
				Y:     dt.Double(userPosition.Y),
				Z:     dt.Double(userPosition.Z),
				Yaw:   dt.Float(userPosition.Yaw),
				Pitch: dt.Float(userPosition.Pitch),
				Flags: 0b0000_0000,
			})

			packets = append(packets, s2c.SetDefaultSpawnPosition{
				Location: dt.Position{X: int(userPosition.X), Y: int(userPosition.Y), Z: int(userPosition.Z)},
			})

			packets = append(packets, s2c.SystemChatMessage{
				Content:         fmt.Sprintf("[+] %s joined the server.", loginSuccess.Username),
				Overlay:         false,
				ShouldBroadcast: true,
			})

			// FIXME Insert both actions in the same PlayerInfoUpdate packet
			packets = append(packets, s2c.PlayerInfoUpdate{
				UUID: loginSuccess.UUID,
				ActionList: []s2c.Action{
					s2c.ActionAddPlayer{
						Name:       loginSuccess.Username,
						Properties: []dt.Property{},
					},
				},

				ShouldBroadcast: true,
			})

			packets = append(packets, s2c.PlayerInfoUpdate{
				UUID: loginSuccess.UUID,
				ActionList: []s2c.Action{
					s2c.ActionUpdateListed{
						Listed: true,
					},
				},

				ShouldBroadcast: true,
			})

			packets = append(packets, s2c.PlaySpawnPlayer{
				EntityId:        dt.VarInt(entityId),
				UUID:            loginSuccess.UUID,
				X:               dt.Double(userPosition.X),
				Y:               dt.Double(userPosition.Y),
				Z:               dt.Double(userPosition.Z),
				Yaw:             dt.UByte(userPosition.Yaw),
				Pitch:           dt.UByte(userPosition.Pitch),
				ShouldBroadcast: true,
			})

		}
	}

	return packets, nil
}

func onSuccessfulLogin() []Packet {
	packets := []Packet{
		s2c.LoginPlay{},
	}

	// send 7X7 chunk square
	for _, chunkPacket := range managers.GetBlockManager().GetChunks() {
		packets = append(packets, chunkPacket)
	}

	return packets
}
