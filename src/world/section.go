package world

import "mango/src/world/block"

type Section struct {
	Blocks [16][16][16]block.Block // [y][z][x]
}

var EMPTY_SECTION Section = GenerateSection(block.AIR)

func GenerateSection(b block.Block) Section {
	var blocks [16][16][16]block.Block

	palette := make([]block.Block, 256)
	for i := 0; i < len(palette); i++ {
		palette[i] = block.Block{ID: i + 1, Name: ""}
	}

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			for k := 0; k < 16; k++ {
				if b.ID == 0 {
					blocks[i][j][k] = b
				} else {
					blocks[i][j][k] = palette[(i+j+k)%len(palette)]
				}
			}
		}
	}

	return Section{
		Blocks: blocks,
	}
}
