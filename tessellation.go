package gen

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/fogleman/delaunay"
	"math"
	"math/rand"
)

// triangulate creates a delaunay.Triangulation centred around the world.ChunkPos passed. The triangulation created is
// deterministic for that specific world.ChunkPos.
// The diagramRadius passed specifies the radius of the triangulation. A bigger radius means more accurately aligning cells
// are created in neighbouring chunks. Depending on the pointDensity passed, which indirectly influences the density of
// voronoi cells, the diagramRadius should be increased or decreased. A higher pointDensity results in smaller cells,
// which doesn't require as large of a triangulation to be accurate. The pointDensity is a value [0-1) that specifies the
// likelihood that a chunk contains one node of a voronoi cell.
func triangulate(pos world.ChunkPos, diagramRadius int32, pointDensity float64) *delaunay.Triangulation {
	r := rand.New(rand.NewSource(chunkHash(pos)))
	d := float64(diagramRadius*2 + 1)

	// Make a rough estimate of the amount of points we'll generate. Assuming the chance a point is generated in a chunk
	// is pointDensity/1, we should be able to multiply that by the total amount of chunks and get a rough estimate.
	points := make([]delaunay.Point, 0, int(d*d*pointDensity))

	for x := -diagramRadius; x <= diagramRadius; x++ {
		for z := -diagramRadius; z <= diagramRadius; z++ {
			// Seed the random instance with the chunk hash of this specific chunk position so that we get deterministic
			// results per chunk.
			r.Seed(chunkHash(world.ChunkPos{pos[0] + x, pos[1] + z}))

			// Increase the density based on a random value, this makes it possible to have more detailed and less
			// consistent biome edges in some places and reduces the general consistency of point spacing.
			density := pointDensity * (1 + r.Float64())

			if r.Float64() < density {
				// We need to generate a point: To obtain more random voronoi cells, we add another 0-15 to every
				// produced coordinate.
				v := r.Int31()
				points = append(points, delaunay.Point{
					X: float64(x<<4 + (v & 0xf)),
					Y: float64(z<<4 + ((v >> 4) & 0xf)),
				})
			}
		}
	}
	t, err := delaunay.Triangulate(points)
	if err != nil {
		panic(err)
	}
	return t
}

type cell struct {
	corners []delaunay.Point
}

func (c cell) inside(v delaunay.Point) bool {
	l := len(c.corners) - 1
	for i, a := range c.corners {
		b := c.corners[l]
		if i != 0 {
			b = c.corners[i-1]
		}
		if (b.X-a.X)*(v.Y-a.Y)-(b.Y-a.Y)*(v.X-a.X) > 0 {
			return false
		}
	}
	return true
}

func (c cell) centre() delaunay.Point {
	p := delaunay.Point{}
	for _, corner := range c.corners {
		p.X += corner.X
		p.Y += corner.Y
	}
	l := float64(len(c.corners))
	p.X /= l
	p.Y /= l
	return p
}

func voronoiCells(d *delaunay.Triangulation) []cell {
	pointList := make([]cell, 0, 15)
	m := make(map[int]struct{})
	for e := 0; e < len(d.Triangles); e++ {
		p := d.Triangles[nextHalfEdge(e)]
		if _, ok := m[p]; !ok {
			m[p] = struct{}{}
			triangles := edgesAroundPoint(d, e)
			for i, v := range triangles {
				triangles[i] = triangleOfEdge(v)
			}
			points := make([]delaunay.Point, len(triangles))
			for i, t := range triangles {
				points[i] = triangleCentre(d, t)
			}
			pointList = append(pointList, cell{corners: points})
		}
	}
	return pointList
}

func edgesAroundPoint(d *delaunay.Triangulation, start int) []int {
	current := start
	p := make([]int, 0, 10)
	for {
		p = append(p, current)
		current = d.Halfedges[nextHalfEdge(current)]
		if current == start || current == -1 {
			break
		}
	}
	return p
}

func iterateVoronoiEdges(d *delaunay.Triangulation, f func(pos delaunay.Point)) {
	for e := 0; e < len(d.Triangles); e++ {
		if e < d.Halfedges[e] {
			a := triangleCentre(d, triangleOfEdge(e))
			b := triangleCentre(d, triangleOfEdge(d.Halfedges[e]))

			diffX, diffY := b.X-a.X, b.Y-a.Y
			dist := math.Sqrt(diffX*diffX + diffY*diffY)
			stepX, stepY := diffX/dist, diffY/dist

			for step := 0.0; step < dist; step++ {
				f(delaunay.Point{X: a.X + stepX*step, Y: a.Y + stepY*step})
			}
		}
	}
}

func triangleCentre(d *delaunay.Triangulation, t int) delaunay.Point {
	points := pointsOfTriangle(d, t)
	return circumcenter(points[0], points[1], points[2])
}

func circumcenter(a, b, c delaunay.Point) delaunay.Point {
	ad := a.X*a.X + a.Y*a.Y
	bd := b.X*b.X + b.Y*b.Y
	cd := c.X*c.X + c.Y*c.Y
	D := 2 * (a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y))
	return delaunay.Point{
		X: 1 / D * (ad*(b.Y-c.Y) + bd*(c.Y-a.Y) + cd*(a.Y-b.Y)),
		Y: 1 / D * (ad*(c.X-b.X) + bd*(a.X-c.X) + cd*(b.X-a.X)),
	}
}

func edgesOfTriangle(t int) [3]int {
	return [3]int{3 * t, 3*t + 1, 3*t + 2}
}

func pointsOfTriangle(d *delaunay.Triangulation, t int) [3]delaunay.Point {
	e := edgesOfTriangle(t)
	return [3]delaunay.Point{d.Points[d.Triangles[e[0]]], d.Points[d.Triangles[e[1]]], d.Points[d.Triangles[e[2]]]}
}

func triangleOfEdge(e int) int {
	return int(math.Floor(float64(e) / 3))
}

func nextHalfEdge(e int) int {
	if e%3 == 2 {
		return e - 2
	}
	return e + 1
}

// chunkHash produces a unique chunk hash for the position passed for chunk positions up to math.MaxInt16. Go random
// seeds only use 32 bits of the int64 you pass, so this function only uses the first 32 bits of the int64.
func chunkHash(pos world.ChunkPos) int64 {
	return int64((uint64(uint16(pos[0])) << 16) | uint64(uint16(pos[1])))
}
