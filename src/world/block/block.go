package block

type Block struct {
	ID   int
	Name string
}

var (
	AIR      = Block{0, "minecraft:air"}
	STONE    = Block{1, "minecraft:stone"}
	ANDESITE = Block{0, "minecraft:andesite"}
	GRANITE  = Block{0, "minecraft:granite"}
	DIORITE  = Block{2, "minecraft:diorite"}
	BEDROCK  = Block{0, "minecraft:bedrock"}
	DIRT     = Block{0, "minecraft:dirt"}
	GRASS    = Block{0, "minecraft:grass"}
)
