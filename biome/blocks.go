package biome

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	grass = world.BlockRuntimeID(block.Grass{})
	dirt  = world.BlockRuntimeID(block.Dirt{})
	stone = world.BlockRuntimeID(block.Stone{})
	sand  = world.BlockRuntimeID(block.Sand{})
)
