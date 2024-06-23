package block

type Block struct {
	ID   int
	Name string
}

var (
	// FIXME: These are not the real IDs
	AIR      = Block{0, "minecraft:air"}
	STONE    = Block{1, "minecraft:stone"}
	ANDESITE = Block{2, "minecraft:andesite"}
	GRANITE  = Block{3, "minecraft:granite"}
	DIORITE  = Block{4, "minecraft:diorite"}
	BEDROCK  = Block{5, "minecraft:bedrock"}
	DIRT     = Block{6, "minecraft:dirt"}
	GRASS    = Block{7, "minecraft:grass"}
)
