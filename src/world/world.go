package world

import (
	"mango/src/world/block"
)

var WorldInstance = NewWorld()

type World struct {
	Chunks map[ChunkPos]*Chunk
}

func NewWorld() *World {
	return &World{
		Chunks: make(map[ChunkPos]*Chunk),
	}
}

func (w *World) GetChunk(pos ChunkPos) *Chunk {
	chunk, exists := w.Chunks[pos]

	if !exists {
		chunk = w.generateOrLoadChunk(pos)
		w.Chunks[pos] = chunk
	}
	return chunk
}

func (w *World) SetChunk(chunk *Chunk) {
	w.Chunks[chunk.Position] = chunk
}

func (w *World) GetBlock(pos BlockPos) block.Block {
	chunk := w.GetChunk(pos.ToChunkPos())
	return chunk.GetBlockAt(pos)
}

func (w *World) SetBlock(pos BlockPos, block block.Block) {
	chunk := w.GetChunk(pos.ToChunkPos())
	chunk.SetBlockAt(pos, block)
}

func (w *World) RemoveBlock(pos BlockPos) {
	w.SetBlock(pos, block.AIR)
}

func (w *World) generateOrLoadChunk(pos ChunkPos) *Chunk {
	return GenerateChunk(pos)
}
