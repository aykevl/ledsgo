package ledsgo

import (
	"math"
	"math/rand"
	"testing"
)

func TestSqrt(t *testing.T) {
	// Test all positive 32-bit integers with the following loop:
	//     for i := (1 << 31) - 1; i >= 0; i-- {
	// Test a large number of positive 32-bit integers.
	for i := 0; i < 100000; i++ {
		n := int(rand.Int31())
		n1 := Sqrt(n)
		n2 := int(math.Sqrt(float64(n)))
		diff := n1 - n2
		if diff < 0 {
			diff = -diff
		}
		if diff > 1 {
			t.Fatalf("sqrt failed: i=%d n1=%d n2=%d diff=%d\n", i, n1, n2, diff)
		}
	}
}

func TestBlend(t *testing.T) {
	// Test a few common cases that must always be correct.
	for _, tc := range []struct {
		bottom uint8
		top    uint8
		alpha  uint8
		result uint8
	}{
		// 0..255
		{0, 255, 0, 0},
		{0, 255, 255, 255},
		{0, 255, 15, 15},
		{0, 255, 240, 240},

		// 255..0
		{255, 0, 0, 255},
		{255, 0, 255, 0},
		{255, 0, 15, 240},
		{255, 0, 240, 15},
	} {
		testBlend(t, tc.bottom, tc.top, tc.alpha, tc.result)
	}
	if t.Failed() {
		return // don't clutter the test output
	}

	// Test the 0..255 scale exhaustively.
	for i := 0; i <= 255; i++ {
		testBlend(t, 0, 255, uint8(i), uint8(i))
		testBlend(t, 255, 0, uint8(i), uint8(255-i))
	}

	// Compare test results against ideally rounded value (using floats).
	// TODO: improve blending so it does a better job at rounding. Now, it can
	// be slightly off due to integer division.
	diffSum := 0.0
	const tests = 1e6
	for i := 0; i < tests; i++ {
		top := uint8(rand.Uint32())
		bottom := uint8(rand.Uint32())
		alpha := uint8(rand.Uint32())
		floatAlpha := float64(alpha) / 255
		ideal := float64(top)*floatAlpha + float64(bottom)*(1-floatAlpha)
		actual := blend(bottom, top, alpha)
		diff := float64(actual) - ideal
		diffSum += math.Abs(diff)
		if math.Abs(diff) >= 1.0 {
			t.Errorf("top %3d  bottom %3d  alpha %3d:  got %3d  expected %.1f (diff: %.3f)", top, bottom, alpha, actual, ideal, diff)
		}
	}
	t.Logf("diff avg: %.2f", diffSum/tests)
}

func testBlend(t *testing.T, bottom, top, alpha, expected uint8) {
	actual := blend(bottom, top, alpha)
	if expected != actual {
		t.Errorf("bottom %3d  top %3d  alpha %3d:  expected %3d, got %d", bottom, top, alpha, expected, actual)
	}
}
