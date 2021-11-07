package gen

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"math"
)

// terrainColumn is a column part of a terrainMap. It holds information on the Biome in that column and the height
// produced by that biome.
type terrainColumn struct {
	height float64
	biome  Biome
}

// terrainMap holds terrain information about an area of the world in the form of columns with heights and biomes.
type terrainMap []terrainColumn

// calculateTerrainMap calculates a terrainMap at a specific world.ChunkPos. The r value passed specifies how much space
// around the chunk's bounds is also calculated to prepare for smoothing the terrain map.
func calculateTerrainMap(r int, pos world.ChunkPos, g *Generator, chunk *chunk.Chunk) terrainMap {
	d := triangulate(pos, 20, 0.06)
	cells := voronoiCells(d)
	g.displayDiagram(d, 128, chunk)

	baseX, baseY := int(pos[0]<<4), int(pos[1]<<4)

	dx := 2*r + 16
	m := make(terrainMap, dx*dx)

	for x := -r; x < 16+r; x++ {
		for y := -r; y < 16+r; y++ {
			biome := g.biome(int32(baseX), int32(baseY), int32(x+baseX), int32(y+baseY), cells)

			m[(x+r)+(y+r)*dx] = terrainColumn{
				height: biome.Height(float64(baseX+x), float64(baseY+y)) * 128,
				biome:  biome,
			}
		}
	}
	return m
}

// smooth smooths the terrainMap where r specifies the radius of the circle around a column that influences the final
// height of a block. The curve passed has an influence on the weight of another height around a column at a specific
// distance.
func (m terrainMap) smooth(r int, c curve) terrainMap {
	var (
		rf            = float64(r)
		curveStepSize = float64(len(c)) / rf
		dx            = 16
		smooth        = make(terrainMap, dx*dx)
		norm, height  float64
		biome         Biome
	)

	for x := 0; x < dx; x++ {
		for y := 0; y < dx; y++ {
			norm, height = 0, 0
			thisCol := m[(x+r)+(y+r)*(dx+r*2)]
			thisHeight := thisCol.height
			biome = thisCol.biome

			for xx := -r; xx <= r; xx++ {
				for yy := -r; yy <= r; yy++ {
					if xx == 0 && yy == 0 {
						continue
					}
					dist := math.Sqrt(float64(xx*xx) + float64(yy*yy))
					if dist > rf {
						// The block fell outside of the circle so we don't need to check this. These blocks have a relatively
						// small radius and ignoring them makes for an improvement in performance without losing accuracy.
						continue
					}
					weight := c.at(int(math.Floor(curveStepSize * dist)))
					norm += weight
					col := m[(x+xx+r)+(y+yy+r)*(dx+r*2)]
					h := col.height
					if col.biome == biome {
						h = thisHeight
					}
					height += weight * h
				}
			}
			smooth[x+y*dx] = terrainColumn{
				height: height / norm,
				biome:  biome,
			}
		}
	}
	return smooth
}
