package mappers

import (
	dt "mango/src/network/datatypes"
	"mango/src/network/packet/s2c"
	"mango/src/world"
	"mango/src/world/block"
)

func MapChunkToPacket(chunk *world.Chunk) s2c.ChunkDataAndLight {
	var packet s2c.ChunkDataAndLight
	packet.ChunkPosition.X = dt.Int(chunk.Position.X)
	packet.ChunkPosition.Z = dt.Int(chunk.Position.Z)

	packet.ChunkSections = make([]dt.ChunkSection, 24)
	for i, section := range chunk.Sections {
		packet.ChunkSections[i] = mapChunkSection(*section)
	}

	return packet
}

func mapChunkSection(section world.Section) dt.ChunkSection {
	var chunkSection dt.ChunkSection

	paletteEntries := make(map[block.Block]struct{})
	nonAirBlocks := 0
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			for k := 0; k < 16; k++ {
				b := section.Blocks[i][j][k]
				paletteEntries[b] = struct{}{}
				if b != block.AIR {
					nonAirBlocks++
				}
			}
		}
	}
	paletteLen := len(paletteEntries)

	var BPE uint8
	switch {
	case paletteLen == 1: // single
		BPE = 0
	//case paletteLen <= 16: // indirect
	//	BPE = 4
	//case paletteLen <= 32:
	//	BPE = 5
	//case paletteLen <= 64:
	//	BPE = 6
	//case paletteLen <= 128:
	//	BPE = 7
	//case paletteLen <= 256:
	//	BPE = 8
	default: // direct
		BPE = 15
	}

	chunkSection.NonAirBlocks = dt.Short(nonAirBlocks)
	chunkSection.Biomes = dt.PalettedContainer{
		BitsPerEntry: 0,
		Palette:      []dt.VarInt{55}, // the_void
		Data:         []dt.Long{},
	}

	palette := make([]dt.VarInt, paletteLen)
	i := 0
	for entry := range paletteEntries {
		palette[i] = dt.VarInt(entry.ID)
		i++
	}

	chunkSection.BlockStates = dt.PalettedContainer{
		BitsPerEntry: dt.UByte(BPE),
		Palette:      palette,
		Data:         buildDataArray(section, int(BPE), palette),
	}

	return chunkSection
}

func buildDataArray(section world.Section, bitsPerEntry int, palette []dt.VarInt) []dt.Long {
	if bitsPerEntry == 0 {
		return []dt.Long{}
	}

	// ID MAP:  map[ entryID ] positionInPalette

	idMap := make(map[dt.VarInt]dt.VarInt)
	directPalette := bitsPerEntry == 15

	for i, entryID := range palette {
		if directPalette {
			idMap[entryID] = entryID
		} else {
			idMap[entryID] = dt.VarInt(i)
		}
	}

	// compute output size
	entriesPerLong := 64 / bitsPerEntry
	outputSize := (16*16*16 + entriesPerLong - 1) / entriesPerLong
	output := make([]dt.Long, outputSize)

	i := 0
	offset := 0

	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				// map blockID to palette ID
				id := section.Blocks[y][z][x].ID
				value := int(idMap[dt.VarInt(id)])

				// insert bits into long
				bits := extractBits(value, bitsPerEntry)
				output[i] |= dt.Long(bits << offset)

				// update offset and current long index
				offset += bitsPerEntry
				if offset+bitsPerEntry > 64 {
					offset = 0
					i++
				}
			}
		}
	}

	return output
}

func extractBits(value, bits int) int64 {
	return int64(value & ((1 << bits) - 1))
}
