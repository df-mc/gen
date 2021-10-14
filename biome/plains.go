package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Plains struct{}

func (p *Plains) CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk) {
	c.SetRuntimeID(x, int16(height), z, 0, grass)
	c.SetRuntimeID(x, int16(height-1), z, 0, dirt)
	c.SetRuntimeID(x, int16(height-2), z, 0, dirt)
}

func (p *Plains) ModNoise(v float64) float64 {
	return v*0.075 + 0.09
}
