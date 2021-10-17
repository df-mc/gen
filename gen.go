package gen

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/gen/f"
	"github.com/fogleman/delaunay"
	"github.com/ojrac/opensimplex-go"
	"time"
)

type Generator struct {
	n                      f.F
	wave                   opensimplex.Noise
	temp, hum              f.F
	biomeBlurX, biomeBlurZ f.F
	b                      biomeSet
}

// New creates a new Generator that implements world.Generator.
func New() *Generator {
	seed := time.Now().Unix()
	n := noise(seed, 3, 2, 0.5)
	return &Generator{
		n:          n,
		wave:       opensimplex.NewNormalized(seed),
		biomeBlurX: noise(seed+0x00f, 4, 2, 0.5),
		biomeBlurZ: noise(seed+0x0ff, 4, 2, 0.5),
		temp:       noise(seed+0x0f0, 1, 2, 1),
		hum:        noise(seed+0xf00, 1, 2, 1),
		b:          newBiomeSet(seed),
	}
}

// GenerateChunk generates a chunk.Chunk at a world.ChunkPos in the world.
func (c *Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	m := calculateTerrainMap(7, pos, c, chunk).smooth(7, normalCurve)

	baseX, baseZ := pos[0]<<4, pos[1]<<4
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			col := m[x+z*16]

			for y := int16(0); y <= int16(col.height); y++ {
				chunk.SetRuntimeID(x, y, z, 0, stone)
			}
			col.biome.CoverGround(x, z, baseX+int32(x), baseZ+int32(z), int(col.height), chunk)
		}
	}
}

var stone, _ = world.BlockRuntimeID(block.Stone{})

// displayDiagram displays the voronoi.Diagram passed in the chunk.Chunk passed by drawing the edges in the sky.
func (c *Generator) displayDiagram(d *delaunay.Triangulation, height int, chunk *chunk.Chunk) {
	iterateVoronoiEdges(d, func(pos delaunay.Point) {
		if pos.X >= 0 && pos.X < 16 && pos.Y >= 0 && pos.Y < 16 {
			chunk.SetRuntimeID(uint8(pos.X), int16(height), uint8(pos.Y), 0, stone)
		}
	})
}

// cell finds the voronoi.Cell that a position was in in the voronoi.Diagram passed. If the cell was not found, false
// is returned.
func (c *Generator) cell(baseX, baseZ, absX, absZ int32, cells []cell) (cell, bool) {
	v := delaunay.Point{
		X: float64(absX-baseX) + (c.biomeBlurX(float64(absX)*0.9, float64(absZ)*0.9)-0.5)*50,
		Y: float64(absZ-baseZ) + (c.biomeBlurZ(float64(absX)*0.9, float64(absZ)*0.9)-0.5)*50,
	}
	// Search all cells for one that contains the position we've got.
	for _, c := range cells {
		if c.inside(v) {
			return c, true
		}
	}
	return cell{}, false
}

// biome returns the Biome that a position was in based on the cell the position was in in the voronoi.Diagram passed.
func (c *Generator) biome(baseX, baseZ, absX, absZ int32, cells []cell) Biome {
	const freq = 0.05
	ce, ok := c.cell(baseX, baseZ, absX, absZ, cells)
	if !ok {
		// This really never should happen. If no cell is found, it means the settings passed to create the
		// delaunay.Triangulation were not valid.
		panic(fmt.Sprintf("Didn't find biome at [%v, %v]. Increase triangulation radius or increase point density", baseX, baseZ))
	}
	centre := ce.centre()

	biomeX, biomeZ := centre.X+float64(baseX), centre.Y+float64(baseZ)
	return c.b.selectBiome(c.hum(biomeX*freq, biomeZ*freq), c.temp(biomeX*freq, biomeZ*freq))
}
