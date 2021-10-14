package biome

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	grass, _ = world.BlockRuntimeID(block.Grass{})
	dirt, _  = world.BlockRuntimeID(block.Dirt{})
	stone, _ = world.BlockRuntimeID(block.Stone{})
	sand, _  = world.BlockRuntimeID(block.Sand{})
)
