package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Ocean struct {
	Noise func(x, z float64) float64
}

func (o *Ocean) CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk) {
	c.SetRuntimeID(x, int16(height), z, 0, sand)
	c.SetRuntimeID(x, int16(height-1), z, 0, sand)
	c.SetRuntimeID(x, int16(height-2), z, 0, sand)
}

func (o *Ocean) Height(x, z float64) float64 {
	return o.Noise(x, z)*0.05 + 0.025
}
