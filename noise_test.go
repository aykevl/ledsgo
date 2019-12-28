package ledsgo

import (
	"math"
	"math/rand"
	"testing"

	"github.com/larspensjo/Go-simplex-noise/simplexnoise"
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
		n1 := simplexnoise.Noise1(float64(x) / 0x1000)
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
			n1 := simplexnoise.Noise2(float64(x)/0x1000, float64(y)/0x1000)
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
				n1 := simplexnoise.Noise3(float64(x)/0x1000, float64(y)/0x1000, float64(z)/0x1000)
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

func TestNoise4(t *testing.T) {
	r := rand.NewSource(0)
	numTestsSub := 48 // ~2s
	numTests := numTestsSub * numTestsSub * numTestsSub * numTestsSub
	rangemax := 0.0
	rangemin := 0.0
	rangesum := 0.0
	diffsum := 0.0
	diffmax := 0.0
	diffmin := 0.0
	for i := 0; i < numTestsSub; i++ { // .12
		for j := 0; j < numTestsSub; j++ { // .12
			for k := 0; k < numTestsSub; k++ { // .12
				for l := 0; l < numTestsSub; l++ { // .12
					x := int32(r.Int63())
					y := int32(r.Int63())
					z := int32(r.Int63())
					w := int32(r.Int63())
					n1 := simplexnoise.Noise4(float64(x)/0x1000, float64(y)/0x1000, float64(z)/0x1000, float64(w)/0x1000)
					n2 := float64(Noise4(int32(x), int32(y), int32(z), int32(w))) / 0x8000
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
							t.Logf("diffmax at x=%d; y=%d; z=%d; w=%d: %f", x, y, z, w, diffmax)
							t.Fail()
						}
					}
					if diff < diffmin {
						diffmin = diff
						if diffmin < -0.1 {
							t.Logf("diffmin at x=%d; y=%d; z=%d; w=%d: %f", x, y, z, w, diffmin)
							t.Fail()
						}
					}
				}
			}
		}
	}
	rangeavg := rangesum / float64(numTests)
	diffavg := diffsum / float64(numTests)
	t.Logf("number of tests: %d", numTests)
	t.Logf("range: avg %+2.6f max %+2.6f min %+2.6f", rangeavg, rangemax, rangemin)
	t.Logf("diff:  avg %+2.6f max %+2.6f min %+2.6f", diffavg, diffmax, diffmin)
	if diffavg >= 0.0004 {
		t.Errorf("diff avg between float and fixed-point is too big: %f", diffavg)
	}
	if diffmax > 0.005 {
		t.Errorf("diff max is too high: %f", diffmax)
	}
	if diffmin < -0.005 {
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
		r = simplexnoise.Noise1(float64(n))
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
		r = simplexnoise.Noise2(float64(n), float64(n))
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
		r = simplexnoise.Noise3(float64(n), float64(n), float64(n))
	}
	resultFloat64 = r
}

func BenchmarkNoise4(b *testing.B) {
	var r int16
	for n := 0; n < b.N; n++ {
		r = Noise4(int32(n), int32(n), int32(n), int32(n))
	}
	resultInt16 = r
}

func BenchmarkNoise4Float(b *testing.B) {
	var r float64
	for n := 0; n < b.N; n++ {
		r = simplexnoise.Noise4(float64(n), float64(n), float64(n), float64(n))
	}
	resultFloat64 = r
}
