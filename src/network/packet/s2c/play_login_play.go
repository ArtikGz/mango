package s2c

import (
	"bytes"
	"io"
	dt "mango/src/network/datatypes"
	"os"
)

type LoginPlay struct {
	EntityID            dt.Int
	IsHardcore          dt.Boolean
	Gamemode            dt.UByte
	PreviousGamemode    dt.Byte
	DimensionNames      []dt.String // Identifier
	RegistryCodec       dt.NbtCompound
	DimensionType       dt.String // Identifier
	DimensionName       dt.String // Identifier
	HashedSeed          dt.Long
	MaxPlayers          dt.VarInt
	ViewDistance        dt.VarInt
	SimulationDistance  dt.VarInt
	ReducedDebugInfo    dt.Boolean
	EnableRespawnScreen dt.Boolean
	IsDebug             dt.Boolean
	IsFlat              dt.Boolean
	HasDeathLocation    dt.Boolean
	DeathDimensionName  dt.String // Identifier
	DeathLocation       dt.Position
}

func (pk LoginPlay) Bytes() []byte {
	// return pk.getRegularBytes() // TODO: use this version at some point
	return pk.getStoredPacketBytes()
}

// loads the full packet bytes
func (pk LoginPlay) getStoredPacketBytes() []byte {
	f, err := os.Open("bin/fullLoginPacket.bin")
	arr, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return arr
}

// Loads only the registryCodec bytes
func (pk LoginPlay) getStoredRegistryBytes() []byte {
	f, err := os.Open("bin/registryCodec.bin")
	arr, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	n := len(arr)
	return arr[:n]
}

// TODO: use this version at some point
func (pk *LoginPlay) getRegularBytes() []byte {
	pk.populatePacket()

	dimBytes := dt.VarInt(len(pk.DimensionNames)).Bytes()
	for _, dim := range pk.DimensionNames {
		dimBytes = append(dimBytes, dim.Bytes()...)
	}

	var deathLocationBytes []byte
	if pk.HasDeathLocation {
		deathLocationBytes = append(pk.DeathDimensionName.Bytes(), pk.DeathLocation.Bytes()...)
	}

	return bytes.Join([][]byte{
		dt.VarInt(0x28).Bytes(),

		pk.EntityID.Bytes(),
		pk.IsHardcore.Bytes(),
		pk.Gamemode.Bytes(),
		pk.PreviousGamemode.Bytes(),

		dimBytes,

		pk.getStoredRegistryBytes(),

		pk.DimensionType.Bytes(),
		pk.DimensionName.Bytes(),

		pk.HashedSeed.Bytes(),
		pk.MaxPlayers.Bytes(),

		pk.ViewDistance.Bytes(),
		pk.SimulationDistance.Bytes(),

		pk.ReducedDebugInfo.Bytes(),
		pk.EnableRespawnScreen.Bytes(),
		pk.IsDebug.Bytes(),
		pk.IsFlat.Bytes(),
		pk.HasDeathLocation.Bytes(),

		deathLocationBytes,
	}, nil)
}

func (pk *LoginPlay) populatePacket() {
	pk.EntityID = 1
	pk.IsHardcore = false
	pk.Gamemode = 1
	pk.PreviousGamemode = 0xFF // -1 aka undefined
	pk.DimensionNames = []dt.String{
		"minecraft:overworld",
		"minecraft:the_end",
		"minecraft:the_nether",
	}

	pk.RegistryCodec = dt.GetDemoRegistryCodec()

	pk.DimensionType = "minecraft:overworld"
	pk.DimensionName = "minecraft:overworld"

	pk.HashedSeed = 0xEDD1337DeadFace

	pk.MaxPlayers = 69
	pk.ViewDistance = 10
	pk.SimulationDistance = 10

	pk.ReducedDebugInfo = false
	pk.EnableRespawnScreen = true
	pk.IsDebug = false
	pk.IsFlat = true
	pk.HasDeathLocation = false
	pk.DeathDimensionName = ""
	pk.DeathLocation = dt.Position{X: 0, Y: 0, Z: 0}
}

func (pk LoginPlay) Broadcast() bool {
	return false
}
