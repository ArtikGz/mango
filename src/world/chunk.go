package world

import (
	"mango/src/world/block"
)

type Chunk struct {
	Position ChunkPos
	Sections [24]*Section
}

func GenerateChunk(pos ChunkPos) *Chunk {
	var b block.Block
	if (pos.X+pos.Z)%2 == 0 {
		b = block.STONE
	} else {
		b = block.DIORITE
	}

	var sections [24]*Section
	generatedSection := GenerateSection(b)
	sections[0] = &generatedSection

	for i := 1; i < 24; i++ {
		tmp := EMPTY_SECTION
		sections[i] = &tmp
	}

	return &Chunk{Position: pos, Sections: sections}
}

func (c *Chunk) GetSectionAt(pos BlockPos) *Section {
	return c.Sections[8+pos.Y>>4]
}

func (c *Chunk) GetBlockAt(pos BlockPos) block.Block {
	cx, cy, cz := pos.X&15, pos.Y&15, pos.Z&15
	return c.GetSectionAt(pos).Blocks[cy][cx][cz]
}

func (c *Chunk) SetBlockAt(pos BlockPos, block block.Block) {
	cx, cy, cz := pos.X&15, pos.Y&15, pos.Z&15
	c.GetSectionAt(pos).Blocks[cy][cx][cz] = block
}
