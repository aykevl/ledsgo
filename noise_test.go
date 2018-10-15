package ledsgo

// Note: this file contains external code, see below.

import (
	"math"
	"math/rand"
	"testing"
)

func TestNoise1(t *testing.T) {
	numTests := int64(10000000) // ~0.25s
	//numTests = 0xffffffff // for exhaustive testing (takes ~100 seconds)
	rangemax := 0.0
	rangemin := 0.0
	rangesum := 0.0
	diffsum := 0.0
	diffmax := 0.0
	diffmin := 0.0
	for x := int64(0); x < numTests; x++ { // .12
		n1 := Noise1Float(float64(x) / 0x1000)
		n2 := float64(Noise1(int32(x))) / 0x8000
		rangesum += n2
		if n2 > rangemax {
			rangemax = n2
		}
		if n2 < rangemin {
			rangemin = n2
		}
		diff := n1 - n2
		diffsum += math.Abs(diff)
		if diff > diffmax {
			diffmax = diff
		}
		if diff < diffmin {
			diffmin = diff
		}
		//t.Logf("%d: %+2.6f %+2.6f %+2.6f", x, n1, n2, diff)
	}
	rangeavg := rangesum / float64(numTests)
	diffavg := diffsum / float64(numTests)
	// Note: the average output is off by ~0.03, which is expected and similar
	// to the floating point implementation.
	t.Logf("number of tests: %d", numTests)
	t.Logf("range: avg %+2.6f max %+2.6f min %+2.6f", rangeavg, rangemax, rangemin)
	t.Logf("diff:  avg %+2.6f max %+2.6f min %+2.6f", diffavg, diffmax, diffmin)
	if diffavg >= 0.00006 {
		t.Errorf("diff between float and fixed-point is too big: %f", diffavg)
	}
	if diffmax > 0.0003 {
		t.Errorf("max is too high: %f", diffmax)
	}
	if diffmin < -0.0003 {
		t.Errorf("min is too low: %f", diffmin)
	}
}

func TestNoise2(t *testing.T) {
	r := rand.NewSource(0)
	numTestsSub := 4000 // ~2s
	//numTestsSub = 16000 // ~30s
	//numTestsSub = 128000 // ~30m
	numTests := numTestsSub * numTestsSub
	rangemax := 0.0
	rangemin := 0.0
	rangesum := 0.0
	diffsum := 0.0
	diffmax := 0.0
	diffmin := 0.0
	for i := 0; i < numTestsSub; i++ { // .12
		for j := 0; j < numTestsSub; j++ { // .12
			x := int32(r.Int63())
			y := int32(r.Int63())
			n1 := Noise2Float(float64(x)/0x1000, float64(y)/0x1000)
			n2 := float64(Noise2(int32(x), int32(y))) / 0x8000
			rangesum += n2
			if n2 > rangemax {
				rangemax = n2
			}
			if n2 < rangemin {
				rangemin = n2
			}
			diff := n1 - n2
			diffsum += math.Abs(diff)
			if diff > diffmax {
				diffmax = diff
				if diffmax > 0.1 {
					t.Fatalf("diffmax at x=%d; y=%d: %f", x, y, diffmax)
				}
			}
			if diff < diffmin {
				diffmin = diff
				if diffmin < -0.1 {
					t.Fatalf("diffmin at x=%d; y=%d: %f", x, y, diffmin)
				}
			}
			if i%32 == 0 && j%32 == 0 {
				//t.Logf("%11d: %+2.6f %+2.6f %+2.6f", x, n1, n2, diff)
			}
		}
	}
	rangeavg := rangesum / float64(numTests)
	diffavg := diffsum / float64(numTests)
	t.Logf("number of tests: %d", numTests)
	t.Logf("range: avg %+2.6f max %+2.6f min %+2.6f", rangeavg, rangemax, rangemin)
	t.Logf("diff:  avg %+2.6f max %+2.6f min %+2.6f", diffavg, diffmax, diffmin)
	if diffavg >= 0.0008 {
		t.Errorf("diff avg between float and fixed-point is too big: %f", diffavg)
	}
	if diffmax > 0.005 {
		t.Errorf("diff max is too high: %f", diffmax)
	}
	if diffmin < -0.005 {
		t.Errorf("diff min is too low: %f", diffmin)
	}
}

func TestNoise3(t *testing.T) {
	r := rand.NewSource(0)
	numTestsSub := 200 // ~2s
	//numTestsSub = 800 // ~130s
	numTests := numTestsSub * numTestsSub * numTestsSub
	rangemax := 0.0
	rangemin := 0.0
	rangesum := 0.0
	diffsum := 0.0
	diffmax := 0.0
	diffmin := 0.0
	for i := 0; i < numTestsSub; i++ { // .12
		for j := 0; j < numTestsSub; j++ { // .12
			for k := 0; k < numTestsSub; k++ { // .12
				x := int32(r.Int63())
				y := int32(r.Int63())
				z := int32(r.Int63())
				n1 := Noise3Float(float64(x)/0x1000, float64(y)/0x1000, float64(z)/0x1000)
				n2 := float64(Noise3(int32(x), int32(y), int32(z))) / 0x8000
				rangesum += n2
				if n2 > rangemax {
					rangemax = n2
				}
				if n2 < rangemin {
					rangemin = n2
				}
				diff := n1 - n2
				diffsum += math.Abs(diff)
				if diff > diffmax {
					diffmax = diff
					if diffmax > 0.1 {
						t.Fatalf("diffmax at x=%d; y=%d; z=%d: %f", x, y, z, diffmax)
					}
				}
				if diff < diffmin {
					diffmin = diff
					if diffmin < -0.1 {
						t.Fatalf("diffmin at x=%d; y=%d; z=%d: %f", x, y, z, diffmin)
					}
				}
				if i%32 == 0 && j%32 == 0 && k%32 == 0 {
					//t.Logf("%+2.6f %+2.6f %+2.6f", n1, n2, diff)
				}
			}
		}
	}
	rangeavg := rangesum / float64(numTests)
	diffavg := diffsum / float64(numTests)
	t.Logf("number of tests: %d", numTests)
	t.Logf("range: avg %+2.6f max %+2.6f min %+2.6f", rangeavg, rangemax, rangemin)
	t.Logf("diff:  avg %+2.6f max %+2.6f min %+2.6f", diffavg, diffmax, diffmin)
	if diffavg >= 0.0006 {
		t.Errorf("diff avg between float and fixed-point is too big: %f", diffavg)
	}
	if diffmax > 0.008 {
		t.Errorf("diff max is too high: %f", diffmax)
	}
	if diffmin < -0.008 {
		t.Errorf("diff min is too low: %f", diffmin)
	}
}

// avoid compiler optimizations
var (
	resultInt16   int16
	resultFloat64 float64
)

func BenchmarkNoise1(b *testing.B) {
	var r int16
	for n := 0; n < b.N; n++ {
		r = Noise1(int32(n))
	}
	resultInt16 = r
}

func BenchmarkNoise1Float(b *testing.B) {
	var r float64
	for n := 0; n < b.N; n++ {
		r = Noise1Float(float64(n))
	}
	resultFloat64 = r
}

func BenchmarkNoise2(b *testing.B) {
	var r int16
	for n := 0; n < b.N; n++ {
		r = Noise2(int32(n), int32(n))
	}
	resultInt16 = r
}

func BenchmarkNoise2Float(b *testing.B) {
	var r float64
	for n := 0; n < b.N; n++ {
		r = Noise2Float(float64(n), float64(n))
	}
	resultFloat64 = r
}

func BenchmarkNoise3(b *testing.B) {
	var r int16
	for n := 0; n < b.N; n++ {
		r = Noise3(int32(n), int32(n), int32(n))
	}
	resultInt16 = r
}

func BenchmarkNoise3Float(b *testing.B) {
	var r float64
	for n := 0; n < b.N; n++ {
		r = Noise3Float(float64(n), float64(n), float64(n))
	}
	resultFloat64 = r
}

// The following code was written by Stefan Gustavson and was converted to Go by
// Lars Pensjö. It is released in the public domain.
//
// Original note:
//
//     SimplexNoise1234, Simplex noise with true analytic
//     derivative in 1D to 4D.
//
//     Author: Stefan Gustavson, 2003-2005
//     Contact: stegu@itn.liu.se
//
//     This code was GPL licensed until February 2011.
//     As the original author of this code, I hereby
//     release it into the public domain.
//     Please feel free to use it for whatever you want.
//     Credit is appreciated where appropriate, and I also
//     appreciate being told where this code finds any use,
//     but you may do as you like.
//
//
//     C implementation of Perlin Simplex Noise over 1,2,3, and 4 dimensions.
//     Author: Stefan Gustavson (stegu@itn.liu.se)
//
//     Adapted to Go by Lars Pensjö (lars.pensjo@gmail.com)
//     This implementation is "Simplex Noise" as presented by
//     Ken Perlin at a relatively obscure and not often cited course
//     session "Real-Time Shading" at Siggraph 2001 (before real
//     time shading actually took on), under the title "hardware noise".
//     The 3D function is numerically equivalent to his Java reference
//     code available in the PDF course notes, although I re-implemented
//     it from scratch to get more readable code. The 1D, 2D and 4D cases
//     were implemented from scratch by me from Ken Perlin's text.

func FASTFLOOR(x float64) int {
	if x > 0 {
		return int(x)
	}
	return int(x) - 1
}

// Helper functions to compute gradients-dot-residualvectors (1D to 4D)
// Note that these generate gradients of more than unit length. To make
// a close match with the value range of classic Perlin noise, the final
// noise values need to be rescaled to fit nicely within [-1,1].
// (The simplex noise functions as such also have different scaling.)
// Note also that these noise functions are the most practical and useful
// signed version of Perlin noise. To return values according to the
// RenderMan specification from the SL noise() and pnoise() functions,
// the noise values need to be scaled and offset to [0,1], like this:
// float SLnoise = (noise(x,y,z) + 1.0) * 0.5;

func qFloat(cond bool, v1 float64, v2 float64) float64 {
	if cond {
		return v1
	}
	return v2
}

func grad1Float(hash uint8, x float64) float64 {
	h := hash & 15
	grad := float64(1 + h&7) // Gradient value 1.0, 2.0, ..., 8.0
	if h&8 != 0 {
		grad = -grad // Set a random sign for the gradient
	}
	return grad * x // Multiply the gradient with the distance
}

func grad2Float(hash uint8, x float64, y float64) float64 {
	h := hash & 7            // Convert low 3 bits of hash code
	u := qFloat(h < 4, x, y) // into 8 simple gradient directions,
	v := qFloat(h < 4, y, x) // and compute the dot product with (x,y).
	return qFloat(h&1 != 0, -u, u) + qFloat(h&2 != 0, -2*v, 2*v)
}

func grad3Float(hash uint8, x, y, z float64) float64 {
	h := hash & 15                                          // Convert low 4 bits of hash code into 12 simple
	u := qFloat(h < 8, x, y)                                // gradient directions, and compute dot product.
	v := qFloat(h < 4, y, qFloat(h == 12 || h == 14, x, z)) // Fix repeats at h = 12 to 15
	return qFloat(h&1 != 0, -u, u) + qFloat(h&2 != 0, -v, v)
}

// 1D simplex noise
func Noise1Float(x float64) float64 {
	i0 := FASTFLOOR(x)
	i1 := i0 + 1
	x0 := x - float64(i0)
	x1 := x0 - 1

	t0 := 1 - x0*x0
	t0 *= t0
	n0 := t0 * t0 * grad1Float(perm[i0&0xff], x0)

	t1 := 1 - x1*x1
	t1 *= t1
	n1 := t1 * t1 * grad1Float(perm[i1&0xff], x1)
	// The maximum value of this noise is 8*(3/4)^4 = 2.53125
	// A factor of 0.395 would scale to fit exactly within [-1,1].
	// fmt.Printf("Noise1 x %.4f, i0 %v, i1 %v, x0 %.4f, x1 %.4f, perm0 %d, perm1 %d: %.4f,%.4f\n", x, i0, i1, x0, x1, perm[i0&0xff], perm[i1&0xff], n0, n1)
	// The algorithm isn't perfect, as it is assymetric. The correction will normalize the result to the interval [-1,1], but the average will be off by 3%.
	return (n0 + n1 + 0.076368899) / 2.45488110001
}

// 2D simplex noise
func Noise2Float(x, y float64) float64 {
	const F2 = 0.366025403 // F2 = 0.5*(sqrt(3.0)-1.0)
	const G2 = 0.211324865 // G2 = (3.0-Math.sqrt(3.0))/6.0

	var n0, n1, n2 float64 // Noise contributions from the three corners

	// Skew the input space to determine which simplex cell we're in
	s := (x + y) * F2 // Hairy factor for 2D
	xs := x + s
	ys := y + s
	i := FASTFLOOR(xs)
	j := FASTFLOOR(ys)

	t := float64(i+j) * G2
	X0 := float64(i) - t // Unskew the cell origin back to (x,y) space
	Y0 := float64(j) - t
	x0 := x - X0 // The x,y distances from the cell origin
	y0 := y - Y0

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
	var i1, j1 int // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		i1 = 1
		j1 = 0 // lower triangle, XY order: (0,0)->(1,0)->(1,1)
	} else {
		i1 = 0
		j1 = 1
	} // upper triangle, YX order: (0,0)->(0,1)->(1,1)

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6

	x1 := x0 - float64(i1) + G2 // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1 + 2*G2 // Offsets for last corner in (x,y) unskewed coords
	y2 := y0 - 1 + 2*G2

	// Calculate the contribution from the three corners
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad2Float(perm[(i+int(perm[j&0xff]))&0xff], x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad2Float(perm[(i+i1+int(perm[((j+j1)&0xff)&0xff]))&0xff], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad2Float(perm[(i+1+int(perm[(j+1)&0xff]))&0xff], x2, y2)
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to return values in the interval [-1,1].
	return (n0 + n1 + n2) / 0.022108854818853867
}

// 3D simplex noise
func Noise3Float(x, y, z float64) float64 {
	// Simple skewing factors for the 3D case
	const F3 = 0.333333333
	const G3 = 0.166666667

	var n0, n1, n2, n3 float64 // Noise contributions from the four corners

	// Skew the input space to determine which simplex cell we're in
	s := (x + y + z) * F3 // Very nice and simple skew factor for 3D
	xs := x + s
	ys := y + s
	zs := z + s
	i := FASTFLOOR(xs)
	j := FASTFLOOR(ys)
	k := FASTFLOOR(zs)

	t := float64(i+j+k) * G3
	X0 := float64(i) - t // Unskew the cell origin back to (x,y,z) space
	Y0 := float64(j) - t
	Z0 := float64(k) - t
	x0 := float64(x) - X0 // The x,y,z distances from the cell origin
	y0 := float64(y) - Y0
	z0 := float64(z) - Z0

	// For the 3D case, the simplex shape is a slightly irregular tetrahedron.
	// Determine which simplex we are in.
	var i1, j1, k1 int // Offsets for second corner of simplex in (i,j,k) coords
	var i2, j2, k2 int // Offsets for third corner of simplex in (i,j,k) coords

	// This code would benefit from a backport from the GLSL version!
	if x0 >= y0 {
		if y0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0 // X Y Z order
		} else if x0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 0
			k2 = 1 // X Z Y order
		} else {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 1
			j2 = 0
			k2 = 1 // Z X Y order
		}
	} else { // x0<y0
		if y0 < z0 {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 0
			j2 = 1
			k2 = 1 // Z Y X order
		} else if x0 < z0 {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 0
			j2 = 1
			k2 = 1 // Y Z X order
		} else {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0 // Y X Z order
		}
	}

	// A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
	// a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
	// a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
	// c = 1/6.

	x1 := x0 - float64(i1) + G3 // Offsets for second corner in (x,y,z) coords
	y1 := y0 - float64(j1) + G3
	z1 := z0 - float64(k1) + G3
	x2 := x0 - float64(i2) + 2*G3 // Offsets for third corner in (x,y,z) coords
	y2 := y0 - float64(j2) + 2*G3
	z2 := z0 - float64(k2) + 2*G3
	x3 := x0 - 1 + 3*G3 // Offsets for last corner in (x,y,z) coords
	y3 := y0 - 1 + 3*G3
	z3 := z0 - 1 + 3*G3

	// Calculate the contribution from the four corners
	t0 := 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad3Float(perm[(i+int(perm[(j+int(perm[k&0xff]))&0xff]))&0xff], x0, y0, z0)
	}

	t1 := 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad3Float(perm[(i+i1+int(perm[(j+j1+int(perm[(k+k1)&0xff]))&0xff]))&0xff], x1, y1, z1)
	}

	t2 := 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad3Float(perm[(i+i2+int(perm[(j+j2+int(perm[(k+k2)&0xff]))&0xff]))&0xff], x2, y2, z2)
	}

	t3 := 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0
	} else {
		t3 *= t3
		n3 = t3 * t3 * grad3Float(perm[(i+1+int(perm[(j+1+int(perm[(k+1)&0xff]))&0xff]))&0xff], x3, y3, z3)
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to stay just inside [-1,1]
	return (n0 + n1 + n2 + n3) / 0.030555466710745972
}
