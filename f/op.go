package f

import (
	"math"
)

type F func(x, y float64) float64

func Sum(f ...F) F {
	if len(f) == 1 {
		return f[0]
	}
	return func(x, y float64) float64 {
		var i float64
		for _, fu := range f {
			i += fu(x, y)
		}
		return i
	}
}

func Mul(f F, v float64) F {
	return func(x, y float64) float64 {
		return f(x, y) * v
	}
}

func MulF(a, b F) F {
	return func(x, y float64) float64 {
		return a(x, y) * b(x, y)
	}
}

func Deriv(f F, dx float64) F {
	return func(x, y float64) float64 {
		x1 := f(x-dx/2, y)
		x2 := f(x+dx/2, y)

		y1 := f(x, y-dx/2)
		y2 := f(x, y+dx/2)

		dX, dY := (x2-x1)/dx, (y2-y1)/dx

		return math.Sqrt(dX*dX+dY*dY) * 10
	}
}

func WarpDomain(f F, freq, warp float64) F {
	i, j, k, l := 0.0, 0.0, 5.3, 1.3
	return func(x, y float64) float64 {
		warpX, warpY := f(x*freq+i, y*freq+j), f(x*freq+k, y*freq*l)
		return f(x*freq+warpX*warp, y*freq+warpY*warp)
	}
}
