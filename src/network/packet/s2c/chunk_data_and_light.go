package s2c

import (
	"mango/src/nbt"
	dt "mango/src/network/datatypes"
	"mango/src/network/packet"
)

type ChunkDataAndLight struct {
	Header        packet.PacketHeader
	ChunkPosition dt.ChunkPosition
	ChunkSections []dt.ChunkSection
}

func (pk ChunkDataAndLight) Bytes() []byte {
	pk.Header.PacketID = 0x24
	var data []byte

	// chunk coords
	data = append(data, pk.ChunkPosition.Bytes()...)

	// heightmap
	heightmaps := [37]int64{}
	for i := 0; i < len(heightmaps); i++ {
		heightmaps[i] = 0
	}

	heightmapNBT := dt.NbtCompound(dt.NBTCompound(map[string]nbt.NBTTag{
		"MOTION_BLOCKING": dt.NBTLongArray(heightmaps[:]),
	}))
	data = append(data, heightmapNBT.Bytes()...)

	// chunk data (20 subchunks)
	chunkData := []byte{}
	for _, section := range pk.ChunkSections {
		chunkData = append(chunkData, section.Bytes()...)
	}

	chunkDataLength := dt.VarInt(len(chunkData))
	data = append(data, chunkDataLength.Bytes()...)
	data = append(data, chunkData...)

	// block entity
	blockEntityCount := dt.VarInt(0)
	data = append(data, blockEntityCount.Bytes()...)

	// light bitsets (full zeros)
	trustEdges := dt.Boolean(true)
	data = append(data, trustEdges.Bytes()...) // byteBuffer.writeBoolean(this.trustEdges);

	skyMask := dt.BitSet{Data: []dt.Long{}}
	data = append(data, skyMask.Bytes()...) // byteBuffer.writeBitSet(this.skyYMask);

	blockYMask := dt.BitSet{Data: []dt.Long{}}
	data = append(data, blockYMask.Bytes()...) // byteBuffer.writeBitSet(this.blockYMask);

	emptySkyYMask := dt.BitSet{Data: []dt.Long{}}
	data = append(data, emptySkyYMask.Bytes()...) // byteBuffer.writeBitSet(this.emptySkyYMask);

	emptyBlockYMask := dt.BitSet{Data: []dt.Long{}}
	data = append(data, emptyBlockYMask.Bytes()...) // byteBuffer.writeBitSet(this.emptyBlockYMask);

	// collections
	skyUpdatesCount := dt.VarInt(0)
	data = append(data, skyUpdatesCount.Bytes()...)

	// skyUpdatesLength := dt.VarInt(0)
	// data = append(data, skyUpdatesLength.Bytes()...)
	// byteBuffer.writeCollection(this.skyUpdates, FriendlyByteBuf::writeByteArray);

	blockUpdatesCount := dt.VarInt(0)
	data = append(data, blockUpdatesCount.Bytes()...)

	// blockUpdatesLength := dt.VarInt(0)
	// data = append(data, blockUpdatesLength.Bytes()...)
	// byteBuffer.writeCollection(this.blockUpdates, FriendlyByteBuf::writeByteArray);
	//

	pk.Header.WriteHeader(&data)

	return data
}
