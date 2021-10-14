package gen

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/fogleman/delaunay"
	"github.com/ojrac/opensimplex-go"
	"math"
	"time"
)

type Generator struct {
	n                      nf
	wave                   opensimplex.Noise
	temp, hum              nf
	biomeBlurX, biomeBlurZ nf
	b                      biomeList
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
		b:          newBiomeList(n),
	}
}

// GenerateChunk generates a chunk.Chunk at a world.ChunkPos in the world.
func (c *Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	d := triangulate(pos, 18, 0.06)
	cells := voronoiCells(d)

	baseX, baseZ := pos[0]<<4, pos[1]<<4
	for x := int32(0); x < 16; x++ {
		for z := int32(0); z < 16; z++ {
			absX, absZ := baseX+x, baseZ+z

			biome := c.biome(baseX, baseZ, absX, absZ, cells)
			height := int(c.height(baseX, baseZ, absX, absZ, cells))

			ux, uz := uint8(x), uint8(z)
			for y := 0; y < height; y++ {
				chunk.SetRuntimeID(ux, int16(y), uz, 0, stone)
			}
			biome.CoverGround(ux, uz, absX, absZ, height, chunk)
		}
	}

	c.displayDiagram(d, 128, chunk)
}

var stone, _ = world.BlockRuntimeID(block.Stone{})

// height calculates the height at a specific position in the world. It takes the average heights produced by the biomes
// around the position to smooth out the terrain.
func (c *Generator) height(baseX, baseZ, absX, absZ int32, cells []cell) float64 {
	const smoothingRadius = 7
	var normalizer, height float64
	curve := normalCurve

	stepSize := float64(len(curve)) / float64(smoothingRadius)
	noiseVal := c.n(float64(absX), float64(absZ)) + c.wave.Eval2(float64(absX)*0.008, float64(absZ)*0.008)*0.3

	// Iterate through all blocks in the radius to check their biomes.
	for x := -int32(smoothingRadius); x <= smoothingRadius; x++ {
		for z := -int32(smoothingRadius); z <= smoothingRadius; z++ {
			distance := math.Sqrt(float64(x*x) + float64(z*z))
			if distance > smoothingRadius {
				// The block fell outside of the circle so we don't need to check this. These blocks have a relatively
				// small radius and ignoring them makes for an improvement in performance without losing accuracy.
				continue
			}
			index := int(math.Floor(stepSize * distance))
			if index >= len(curve) {
				index = len(curve) - 1
			}

			weight := curve[index]
			biome := c.biome(baseX, baseZ, absX+x, absZ+z, cells)

			normalizer += weight
			height += weight * (128 * biome.ModNoise(noiseVal))
		}
	}
	return height / normalizer
}

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
	localX, localZ := absX-baseX, absZ-baseZ

	v := delaunay.Point{
		X: float64(localX) + (c.biomeBlurX(float64(absX)*0.9, float64(absZ)*0.9)-0.5)*50,
		Y: float64(localZ) + (c.biomeBlurZ(float64(absX)*0.9, float64(absZ)*0.9)-0.5)*50,
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
