package s2c

import (
	"bytes"
	"mango/src/nbt"
	dt "mango/src/network/datatypes"
)

type ChunkDataAndLight struct {
	ChunkPosition dt.ChunkPosition
	ChunkSections []dt.ChunkSection
}

func (pk ChunkDataAndLight) Bytes() []byte {
	// heightmap
	heightmaps := [37]int64{}
	for i := 0; i < len(heightmaps); i++ {
		heightmaps[i] = 0
	}

	heightmapNBT := dt.NbtCompound(dt.NBTCompound(map[string]nbt.NBTTag{
		"MOTION_BLOCKING": dt.NBTLongArray(heightmaps[:]),
	}))

	// chunk data (20 sub chunks)
	chunkData := make([][]byte, len(pk.ChunkSections))
	for i, section := range pk.ChunkSections {
		chunkData[i] = section.Bytes()
	}
	chunkDataBytes := bytes.Join(chunkData, nil)

	// sky mask
	skyMask := dt.BitSet{Data: []dt.Long{}}
	blockYMask := dt.BitSet{Data: []dt.Long{}}
	emptySkyYMask := dt.BitSet{Data: []dt.Long{}}
	emptyBlockYMask := dt.BitSet{Data: []dt.Long{}}

	return bytes.Join([][]byte{
		dt.VarInt(0x24).Bytes(),

		pk.ChunkPosition.Bytes(),

		heightmapNBT.Bytes(),

		dt.VarInt(len(chunkDataBytes)).Bytes(),
		chunkDataBytes,

		// block entity count
		dt.VarInt(0).Bytes(),

		// trust edges
		dt.Boolean(true).Bytes(),

		skyMask.Bytes(),
		blockYMask.Bytes(),
		emptySkyYMask.Bytes(),
		emptyBlockYMask.Bytes(),

		// sky updates count
		dt.VarInt(0).Bytes(),

		// block updates count
		dt.VarInt(0).Bytes(),
	}, nil)
}

func (pk ChunkDataAndLight) Broadcast() bool {
	return false
}
