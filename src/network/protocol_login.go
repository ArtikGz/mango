package network

import (
	"bytes"
	"io"
	"mango/src/config"
	"mango/src/logger"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet"
	"mango/src/network/packet/c2s"
	"mango/src/network/packet/s2c"
)

var worldChunks []s2c.ChunkDataAndLight

func init() {
	worldChunks = generateMockChunks()
}

func generateMockChunks() []s2c.ChunkDataAndLight {
	chunks := make([]s2c.ChunkDataAndLight, 0)

	for i := -3; i < 4; i++ {
		for j := -3; j < 4; j++ {
			var chunkPacket s2c.ChunkDataAndLight

			// Chunk Position
			chunkPacket.ChunkPosition.X = dt.Int(i)
			chunkPacket.ChunkPosition.Z = dt.Int(j)

			// Chunk Sections
			sections := make([]dt.ChunkSection, 0)
			for i := 0; i < 24; i++ {
				nab := 16 * 16 * 16
				if i > 0 {
					nab = 0
				}

				blocks := [16][16][16]dt.Long{}
				for y := 0; y < 16; y++ {
					for z := 0; z < 16; z++ {
						for x := 0; x < 16; x++ {
							blocks[y][z][x] = 1

							if i > 0 {
								blocks[y][z][x] = 0
							}
						}
					}
				}

				section := dt.ChunkSection{
					NonAirBlocks: dt.Short(nab),
					BlockStates: dt.PalettedContainer{
						BitsPerEntry: 8,
						Palette:      []dt.VarInt{0, 1},
						Data:         blocks,
					},
					Biomes: dt.PalettedContainer{
						BitsPerEntry: 0,
						Palette:      []dt.VarInt{55}, // the_void ?
						Data:         [16][16][16]dt.Long{},
					},
				}

				sections = append(sections, section)
			}

			chunkPacket.ChunkSections = sections
			chunks = append(chunks, chunkPacket)
		}
	}

	return chunks
}

func HandleLoginPacket(conn *Connection, data *[]byte) {
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

		} else { // Offline mode, return LoginSuccess
			var logingSuccess s2c.LoginSuccess
			logingSuccess.Username = loginStart.Name
			if loginStart.HasUUID {
				logingSuccess.UUID = loginStart.UUID
			}

			packetBytes := logingSuccess.Bytes()
			conn.outgoingPackets <- &packetBytes
			logger.Debug("Login Success: %+v", logingSuccess)
			conn.state = PLAY

			// send init PLAY packets (Login (Play), Default Spawn Position, etc.)
			onSuccessfulLogin(conn)
		}
	}
}

func onSuccessfulLogin(conn *Connection) {
	var loginPlay s2c.LoginPlay
	packetBytes := loginPlay.Bytes()
	conn.outgoingPackets <- &packetBytes

	var spawnPos s2c.SetDefaultSpawnPosition
	packetBytes1 := spawnPos.Bytes()
	conn.outgoingPackets <- &packetBytes1

	// send 7X7 chunk square
	for _, chunk := range worldChunks {
		packetBytesChunk := chunk.Bytes()
		conn.outgoingPackets <- &packetBytesChunk
	}
}
