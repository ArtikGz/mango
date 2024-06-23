package world

import "mango/src/world/block"

type Section struct {
	Blocks [16][16][16]block.Block // [y][z][x]
}

var EMPTY_SECTION = GenerateSection(block.AIR)

func GenerateSection(b block.Block) Section {
	var blocks [16][16][16]block.Block

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			for k := 0; k < 16; k++ {
				blocks[i][j][k] = b
			}
		}
	}

	return Section{
		Blocks: blocks,
	}
}
