package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Ocean struct{}

func (o *Ocean) CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk) {
	c.SetRuntimeID(x, int16(height), z, 0, sand)
	c.SetRuntimeID(x, int16(height-1), z, 0, sand)
	c.SetRuntimeID(x, int16(height-2), z, 0, sand)
}

func (o *Ocean) ModNoise(v float64) float64 {
	return v*0.05 + 0.025
}
