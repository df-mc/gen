package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/gen/f"
)

type Ocean struct {
	Noise f.F
}

func (o *Ocean) CoverGround(x, z uint8, _, _ int32, height int, c *chunk.Chunk) {
	c.SetBlock(x, int16(height), z, 0, sand)
	c.SetBlock(x, int16(height-1), z, 0, sand)
	c.SetBlock(x, int16(height-2), z, 0, sand)
}

func (o *Ocean) Height(x, z float64) float64 {
	return o.Noise(x, z)*0.05 + 0.025
}
