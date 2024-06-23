package world

type BlockPos struct{ X, Y, Z int }

func (pos BlockPos) ToChunkPos() ChunkPos {
	return ChunkPos{
		X: pos.X >> 4,
		Z: pos.Z >> 4,
	}
}

type ChunkPos struct{ X, Z int }
