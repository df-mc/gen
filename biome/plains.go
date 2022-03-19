package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/gen/f"
)

type Plains struct {
	Noise f.F
}

func (p *Plains) CoverGround(x, z uint8, _, _ int32, height int, c *chunk.Chunk) {
	c.SetBlock(x, int16(height), z, 0, grass)
	c.SetBlock(x, int16(height-1), z, 0, dirt)
	c.SetBlock(x, int16(height-2), z, 0, dirt)
}

func (p *Plains) Height(x, z float64) float64 {
	return p.Noise(x, z)*0.15 + 0.07
}
