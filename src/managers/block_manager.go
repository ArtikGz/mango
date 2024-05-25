package managers

import (
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/s2c"
	"math"
)

var blockManagerInstance blockManager

func init() {
	worldChunks := generateMockChunks()
	blockManagerInstance = blockManager{worldChunks}
}

func generateMockChunks() []s2c.ChunkDataAndLight {
	chunks := make([]s2c.ChunkDataAndLight, 0)

	for cx := -3; cx < 4; cx++ {
		for cz := -3; cz < 4; cz++ {
			var chunkPacket s2c.ChunkDataAndLight

			// Chunk Position
			chunkPacket.ChunkPosition.X = dt.Int(cx)
			chunkPacket.ChunkPosition.Z = dt.Int(cz)

			// Chunk Sections
			sections := make([]dt.ChunkSection, 0)
			for i := 0; i < 24; i++ {
				nab := 0
				if cx == 0 && cz == 0 && i == 0 {
					nab = 16 * 16 * 16
				}

				blocks := [16][16][16]dt.Long{}
				for y := 0; y < 16; y++ {
					for z := 0; z < 16; z++ {
						for x := 0; x < 16; x++ {
							blocks[y][z][x] = 0

							if cx == 0 && cz == 0 && i == 0 {
								blocks[y][z][x] = 1
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

type blockManager struct {
	worldChunks []s2c.ChunkDataAndLight
}

func GetBlockManager() blockManager {
	return blockManagerInstance
}

func (bm blockManager) GetChunks() []s2c.ChunkDataAndLight {
	return bm.worldChunks
}

func calcChunkPos(pos int) int {
	return int(math.Floor(float64(pos) / 16))
}

func (bm blockManager) getChunkAt(x, z int) *s2c.ChunkDataAndLight {
	for i, chunk := range bm.worldChunks {
		if int(chunk.ChunkPosition.X) == x && int(chunk.ChunkPosition.Z) == z {
			return &bm.worldChunks[i]
		}
	}

	return nil
}

func mod(a, b int) int {
	res := int(a % b)
	if (res < 0 && b > 0) || (res > 0 && b < 0) {
		return res + b
	}

	return res
}

func abs(n int) int {
	return int(math.Abs(float64(n)))
}

func (bm blockManager) RemoveBlockAt(x, y, z int) {
	chunkX := calcChunkPos(x)
	chunkZ := calcChunkPos(z)
	chunkY := calcChunkPos(y) + 4

	chunk := bm.getChunkAt(chunkX, chunkZ)

	chunk.ChunkSections[chunkY].NonAirBlocks--
	chunk.ChunkSections[chunkY].BlockStates.Data[mod(y, 16)][mod(z, 16)][mod(x, 16)] = 0
}
