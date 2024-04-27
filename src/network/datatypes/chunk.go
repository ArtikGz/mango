package datatypes

type ChunkPosition struct {
	X Int
	Z Int
}

func (cp *ChunkPosition) Bytes() (buffer []byte) {
	buffer = append(buffer, cp.X.Bytes()...)
	buffer = append(buffer, cp.Z.Bytes()...)
	return
}

type ChunkSection struct {
	NonAirBlocks Short
	BlockStates  PalettedContainer
	Biomes       PalettedContainer
}

func (sc *ChunkSection) Bytes() (buffer []byte) {
	buffer = append(buffer, sc.NonAirBlocks.Bytes()...)
	buffer = append(buffer, sc.BlockStates.Bytes()...)
	buffer = append(buffer, sc.Biomes.Bytes()...)
	return
}

// =====================================================
type PalettedContainer struct {
	BitsPerEntry UByte
	Palette      []VarInt
	Data         [16][16][16]Long
}

func (pc *PalettedContainer) Bytes() (buffer []byte) {
	BPE := int(pc.BitsPerEntry)

	buffer = append(buffer, pc.BitsPerEntry.Bytes()...)

	// single value palette
	if BPE == 0 {
		// TODO: Check if pc.Palette is not empty
		buffer = append(buffer, pc.Palette[0].Bytes()...)

	} else if BPE <= 8 { // indirect palette
		length := VarInt(len(pc.Palette))
		buffer = append(buffer, length.Bytes()...)

		for _, v := range pc.Palette {
			buffer = append(buffer, v.Bytes()...)
		}
	} // else { direct palette (no data) }

	dataLength := VarInt(16 * 16 * 16 * BPE / 64)
	individualValueMask := Long((1 << BPE) - 1)

	data := make([]Long, 0, dataLength)
	for u := 0; u < int(dataLength); u++ {
		data = append(data, 0)
	}

	if dataLength > 0 {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {
					blockNumber := (((y * 16) + z) * 16) + x
					startLong := (blockNumber * BPE) / 64
					startOffset := (blockNumber * BPE) % 64
					endLong := ((blockNumber+1)*BPE - 1) / 64

					block := pc.Data[y][z][x] & individualValueMask

					data[startLong] |= (block << Long(startOffset))

					if startLong != endLong {
						data[endLong] = (block >> (64 - Long(startOffset)))
					}
				}
			}
		}
	}

	buffer = append(buffer, dataLength.Bytes()...)
	for _, v := range data {
		buffer = append(buffer, v.Bytes()...)
	}

	return
}
