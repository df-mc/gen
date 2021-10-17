package biome

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
	"math"
)

type Mountains struct {
	Noise func(x, z float64) float64
}

func (m *Mountains) CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk) {
	s := slope(float64(absX), float64(absZ), m.Height)
	if s < 0.01 {
		c.SetRuntimeID(x, int16(height), z, 0, grass)
	}
}

func (m *Mountains) Height(x, z float64) float64 {
	return m.Noise(x, z) * 0.6
}

// slope calculates roughly the slope at a specific x and z value in the noise function passed.
func slope(x, y float64, noise func(x, z float64) float64) float64 {
	dx, dy := 0.003, 0.003

	x1 := noise(x-dx/2, y)
	x2 := noise(x+dx/2, y)

	y1 := noise(x, y-dy/2)
	y2 := noise(x, y+dy/2)

	dX, dY := (x2-x1)/dx, (y2-y1)/dy

	return math.Sqrt(dX*dX + dY*dY)
}
