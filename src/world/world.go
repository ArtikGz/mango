package world

import (
	"mango/src/logger"
	"mango/src/world/block"
	"sync"
)

var WorldInstance = NewWorld()

type World struct {
	Chunks map[ChunkPos]*Chunk
	mu     sync.RWMutex
}

func NewWorld() *World {
	return &World{
		Chunks: make(map[ChunkPos]*Chunk),
	}
}

func (w *World) GetChunk(pos ChunkPos) *Chunk {
	w.mu.RLock()
	chunk, exists := w.Chunks[pos]
	w.mu.RUnlock()

	if !exists {
		w.mu.Lock()
		defer w.mu.Unlock()

		chunk = w.generateOrLoadChunk(pos)
		w.Chunks[pos] = chunk
	}
	return chunk
}

func (w *World) SetChunk(chunk *Chunk) {
	w.mu.Lock()
	w.Chunks[chunk.Position] = chunk
	w.mu.Unlock()
}

func (w *World) GetBlock(pos BlockPos) block.Block {
	chunk := w.GetChunk(pos.ToChunkPos())
	w.mu.RLock()
	defer w.mu.RUnlock()
	return chunk.GetBlockAt(pos)
}

func (w *World) SetBlock(pos BlockPos, block block.Block) {
	chunk := w.GetChunk(pos.ToChunkPos())
	chunk.SetBlockAt(pos, block)
	w.SetChunk(chunk)
}

func (w *World) RemoveBlock(pos BlockPos) {
	w.SetBlock(pos, block.AIR)
}

func (w *World) generateOrLoadChunk(pos ChunkPos) *Chunk {
	logger.Warn("Generate Chunk %v", pos)
	return GenerateChunk(pos)
}
