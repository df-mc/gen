package gen

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/gen/f"
	"github.com/fogleman/delaunay"
	"time"
)

type Generator struct {
	temp, hum    f.F
	blurX, blurZ f.F
	b            biomeSet
}

// New creates a new Generator that implements world.Generator.
func New() *Generator {
	seed := time.Now().Unix()
	return &Generator{
		blurX: f.Noise(seed+0x00f, 4, 2, 0.5).Norm(),
		blurZ: f.Noise(seed+0x0ff, 4, 2, 0.5).Norm(),
		temp:  f.Noise(seed+0x0f0, 1, 2, 1).Norm(),
		hum:   f.Noise(seed+0xf00, 1, 2, 1).Norm(),
		b:     newBiomeSet(seed),
	}
}

// GenerateChunk generates a chunk.Chunk at a world.ChunkPos in the world.
func (g *Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	m := calculateTerrainMap(10, pos, g, chunk).smooth(10, normalCurve)

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
func (g *Generator) displayDiagram(d *delaunay.Triangulation, height int, chunk *chunk.Chunk) {
	iterateVoronoiEdges(d, func(pos delaunay.Point) {
		if pos.X >= 0 && pos.X < 16 && pos.Y >= 0 && pos.Y < 16 {
			chunk.SetRuntimeID(uint8(pos.X), int16(height), uint8(pos.Y), 0, stone)
		}
	})
}

// cell finds the voronoi.Cell that a position was in in the voronoi.Diagram passed. If the cell was not found, false
// is returned.
func (g *Generator) cell(baseX, baseZ, absX, absZ int32, cells []cell) (cell, bool) {
	v := delaunay.Point{
		X: float64(absX-baseX) + (g.blurX(float64(absX)*0.9, float64(absZ)*0.9)-0.5)*50,
		Y: float64(absZ-baseZ) + (g.blurZ(float64(absX)*0.9, float64(absZ)*0.9)-0.5)*50,
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
func (g *Generator) biome(baseX, baseZ, absX, absZ int32, cells []cell) Biome {
	const freq = 0.05
	ce, ok := g.cell(baseX, baseZ, absX, absZ, cells)
	if !ok {
		// This really never should happen. If no cell is found, it means the settings passed to create the
		// delaunay.Triangulation were not valid.
		panic(fmt.Sprintf("Didn't find biome at [%v, %v]. Increase triangulation radius or increase point density", baseX, baseZ))
	}
	p := ce.centre(float64(baseX), float64(baseZ))
	return g.b.selectBiome(g.hum(p.X*freq, p.Y*freq), g.temp(p.X*freq, p.Y*freq))
}
