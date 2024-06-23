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
	Data         []Long
}

func (pc *PalettedContainer) Bytes() []byte {
	var buffer []byte

	BPE := int(pc.BitsPerEntry)

	// write Bits Per Entry
	buffer = append(buffer, pc.BitsPerEntry.Bytes()...)

	// write Palette
	if BPE == 0 { // single value palette
		buffer = append(buffer, pc.Palette[0].Bytes()...)

	} else if BPE <= 8 { // indirect palette
		buffer = append(buffer, VarInt(len(pc.Palette)).Bytes()...)
		for _, v := range pc.Palette {
			buffer = append(buffer, v.Bytes()...)
		}
	} // else { direct palette (no data) }

	// write Entry list
	buffer = append(buffer, VarInt(len(pc.Data)).Bytes()...)
	for _, entry := range pc.Data {
		buffer = append(buffer, entry.Bytes()...)
	}

	return buffer
}
