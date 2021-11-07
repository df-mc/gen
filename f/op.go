package f

import (
	"math"
)

// F is a function used in the f package. It may be manipulated by any of the methods it has, resulting in a
// new F being created.
type F func(x, y float64) float64

// Sum sums up multiple functions and returns a new F that is the result of adding up the functions.
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

// Norm normalises the results of the old function and returns a new F that returns values in the range [0 1) as
// opposed to (-1 1).
func (f F) Norm() F {
	return func(x, y float64) float64 {
		return (f(x, y) + 1) / 2
	}
}

// Inv inverses the results of the old function and returns a new F that has its values [0 1) changed to [1 0).
func (f F) Inv() F {
	return func(x, y float64) float64 {
		return 1 - f(x, y)
	}
}

// Pow returns a new F that raises values returned by the old function to the power of v.
func (f F) Pow(v float64) F {
	return func(x, y float64) float64 {
		return math.Pow(f(x, y), v)
	}
}

// Thresh returns a new F that returns 1 for values returned by the function that exceed the threshold v. If they
// do not exceed the threshold, 0 is returned.
func (f F) Thresh(v float64) F {
	return func(x, y float64) float64 {
		if f(x, y) > v {
			return 1
		}
		return 0
	}
}

// Abs returns a new F that only returns absolute values.
func (f F) Abs() F {
	return func(x, y float64) float64 {
		return math.Abs(f(x, y))
	}
}

// Mul returns a new F that returns values by the old function multiplied by a value v.
func (f F) Mul(v float64) F {
	return func(x, y float64) float64 {
		return f(x, y) * v
	}
}

// Freq returns a new F with frequency v.
func (f F) Freq(v float64) F {
	return func(x, y float64) float64 {
		return f(x*v, y*v)
	}
}

// MulF returns a new F that returns values of the old function multiplied by the value of the function v at the same
// x and y values.
func (f F) MulF(v F) F {
	return func(x, y float64) float64 {
		return f(x, y) * v(x, y)
	}
}

// Slope returns an F that returns an approximate slope at a position. The dx specified influences the distance over
// which the slope is calculated.
func (f F) Slope(dx float64) F {
	return func(x, y float64) float64 {
		x1 := f(x-dx/2, y)
		x2 := f(x+dx/2, y)

		y1 := f(x, y-dx/2)
		y2 := f(x, y+dx/2)

		dX, dY := (x2-x1)/dx, (y2-y1)/dx

		return math.Sqrt(dX*dX + dY*dY)
	}
}

// WarpDomain returns a new F that performs domain warping on the old function. The frequency passed influences the
// frequency desired of the F returned (practically, this is the zoom level). The warp value specifies the extent to
// which the warping should happen. A value of ~70 generally works well for terrain.
func (f F) WarpDomain(freq, warp float64) F {
	i, j, k, l := 0.0, 0.0, 5.3, 1.3
	return func(x, y float64) float64 {
		warpX, warpY := f(x*freq+i, y*freq+j), f(x*freq+k, y*freq*l)
		return f(x*freq+warpX*warp, y*freq+warpY*warp)
	}
}
