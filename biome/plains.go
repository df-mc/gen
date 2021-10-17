package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Plains struct {
	Noise func(x, z float64) float64
}

func (p *Plains) CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk) {
	c.SetRuntimeID(x, int16(height), z, 0, grass)
	c.SetRuntimeID(x, int16(height-1), z, 0, dirt)
	c.SetRuntimeID(x, int16(height-2), z, 0, dirt)
}

func (p *Plains) Height(x, z float64) float64 {
	return p.Noise(x, z)*0.075 + 0.09
}
