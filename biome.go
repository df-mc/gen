package gen

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/gen/biome"
)

type Biome interface {
	// CoverGround covers the ground of the chunk.Chunk passed using this biome's specific features.
	CoverGround(x, z uint8, absX, absZ int32, height int, c *chunk.Chunk)
	// ModNoise modifies the noise value passed to it according to the landscape features of the biome.
	ModNoise(v float64) float64
}

type biomeList struct {
	Plains    Biome
	Ocean     Biome
	Mountains Biome
}

func newBiomeList(noise nf) biomeList {
	return biomeList{
		Plains:    &biome.Plains{},
		Ocean:     &biome.Ocean{},
		Mountains: &biome.Mountains{Noise: noise},
	}
}

func (b biomeList) selectBiome(hum, temp float64) Biome {
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
