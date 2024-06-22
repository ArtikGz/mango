package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"mango/src/config"
	"mango/src/logger"
	"mango/src/managers"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
	"math/rand/v2"
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

		var uuid []byte
		if loginStart.HasUUID {
			uuid = loginStart.UUID
		} else {
			uuid = make([]byte, 16)
			binary.BigEndian.PutUint64(uuid[:8], rand.Uint64())
			binary.BigEndian.PutUint64(uuid[8:], rand.Uint64())
		}

		user := getOrCreateUser(string(loginStart.Name), uuid)

		if config.IsOnline() {
			// TODO: implement cypher and return EncryptionRequest
			logger.Error("Online mode is not yet supported, please, change online to false in '%s'.", config.GetConfigPath())
		} else { // Offline mode, return LoginSuccess
			if config.CompressionThreshold() > -1 {
				packets = append(packets, s2c.SetCompressionPacket{
					Threshold: dt.VarInt(config.CompressionThreshold()),
				})
			}

			var loginSuccess = s2c.LoginSuccess{
				Username: dt.String(user.Name),
				UUID:     uuid,
			}

			logger.Debug("Login Success: %+v", loginSuccess)

			// send init PLAY packets (Login (Play), Default Spawn Position, etc.)
			packets = append(packets, loginSuccess)
			packets = append(packets, onSuccessfulLogin()...)
			packets = append(packets, []Packet{
				s2c.PlaySynchronizePlayerPosition{
					X:     dt.Double(user.Position.X),
					Y:     dt.Double(user.Position.Y),
					Z:     dt.Double(user.Position.Z),
					Yaw:   dt.Float(user.Position.Yaw),
					Pitch: dt.Float(user.Position.Pitch),
					Flags: 0b0000_0000,
				},
				s2c.SetDefaultSpawnPosition{
					Location: dt.Position{X: int(user.Position.X), Y: int(user.Position.Y), Z: int(user.Position.Z)},
				},
				s2c.SystemChatMessage{
					Content:         fmt.Sprintf("[+] %s joined the server.", loginSuccess.Username),
					Overlay:         false,
					ShouldBroadcast: true,
				},
				// FIXME Insert both actions in the same PlayerInfoUpdate packet
				s2c.PlayerInfoUpdate{
					UUID: loginSuccess.UUID,
					ActionList: []s2c.Action{
						s2c.ActionAddPlayer{
							Name:       loginSuccess.Username,
							Properties: []dt.Property{},
						},
					},

					ShouldBroadcast: true,
				},
				s2c.PlayerInfoUpdate{
					UUID: loginSuccess.UUID,
					ActionList: []s2c.Action{
						s2c.ActionUpdateListed{
							Listed: true,
						},
					},

					ShouldBroadcast: true,
				},
				s2c.PlaySpawnPlayer{
					EntityId:        dt.VarInt(user.EntityId),
					UUID:            loginSuccess.UUID,
					X:               dt.Double(user.Position.X),
					Y:               dt.Double(user.Position.Y),
					Z:               dt.Double(user.Position.Z),
					Yaw:             dt.UByte(user.Position.Yaw),
					Pitch:           dt.UByte(user.Position.Pitch),
					ShouldBroadcast: true,
				},
			}...)
		}
	}

	return packets, nil
}

func getOrCreateUser(name string, uuid []byte) *managers.User {
	user := managers.GetUserManager().GetUser(name)
	if user == nil {
		user = &managers.User{
			EntityId: managers.GetEntityManager().GenerateID(),
			UUID:     uuid,
			Name:     name,
			Position: managers.UserPosition{
				X:     8.5,
				Y:     -46,
				Z:     8.5,
				Yaw:   0,
				Pitch: 0,
			},
		}
		managers.GetUserManager().AddUser(*user)
	}
	return user
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
