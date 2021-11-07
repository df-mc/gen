package f

import (
	"github.com/ojrac/opensimplex-go"
	"math"
)

// Noise returns an function that may be used to calculate a layered noise value at a specific x and y.
// The seed passed is used to seed the noise instances created. The octaves specifies the amount of noise 'layers' that
// should be present. Setting this value to 4 means 4 noise values will be layered on top of each other.
// The lacunarity and persistence specify the frequency and amplitude respectively of subsequent octaves. An example:
// 3 octaves with a lacunarity of 2.0 and a persistence of 0.5 gives us the following noise values layered:
//   0: amp = 0.5^0 = 1.00, freq = 2.0^0 = 1.00
//   1: amp = 0.5^1 = 0.50, freq = 2.0^1 = 2.00
//   2: amp = 0.5^2 = 0.25, freq = 2.0^2 = 4.00
func Noise(seed int64, octaves int, lacunarity, persistence float64) F {
	// Create Noise instances for all octaves.
	n := make([]opensimplex.Noise, octaves)
	for i := 0; i < octaves; i++ {
		n[i] = opensimplex.New(seed)
	}

	// Calculate the maximum possible value when adding together Noise values from all octaves.
	var max float64
	for i := 0; i < octaves; i++ {
		max += math.Pow(persistence, float64(i))
	}

	return func(x, y float64) float64 {
		var v float64

		for i, noise := range n {
			freq := math.Pow(lacunarity, float64(i)) * 0.04
			amp := math.Pow(persistence, float64(i))

			v += amp * noise.Eval2(x*freq, y*freq)
		}
		// Normalise at the end so we get a value in the range [0-1].
		v /= max
		return v
	}
}
