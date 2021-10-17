package gen

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/gen/biome"
)

type Biome interface {
	// CoverGround covers the ground of the chunk.Chunk passed using this biome's specific features.
	CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk)
	// Height returns a height value produced for the biome at a specific x and z in the world. Biomes generally use
	// noise to return a height value.
	Height(x, z float64) float64
}

type biomeSet struct {
	Plains    Biome
	Ocean     Biome
	Mountains Biome
}

func newBiomeSet(seed int64) biomeSet {
	n := noise(seed, 3, 2, 0.5)
	return biomeSet{
		Plains:    &biome.Plains{Noise: n},
		Ocean:     &biome.Ocean{Noise: n},
		Mountains: &biome.Mountains{Noise: n},
	}
}

func (b biomeSet) selectBiome(hum, temp float64) Biome {
	switch {
	case hum < 0.25:
		switch {
		case temp < 0.7:
			return b.Ocean
		case temp < 0.85:
			// river
		default:
			// swamp
		}
	case hum < 0.6:
		switch {
		case temp < 0.25:
			// ice plains
		case temp < 0.75:
			return b.Plains
		default:
			// desert
		}
	case hum < 0.8:
		switch {
		case temp < 0.25:
			// taiga
		case temp < 0.75:
			// forest
		default:
			// birch forest
		}
	default:
		switch {
		case temp < 0.2:
			return b.Mountains
		case temp < 0.4:
			return b.Mountains // small mountains
		default:
			return b.Ocean
		}
	}
	return b.Plains
}
