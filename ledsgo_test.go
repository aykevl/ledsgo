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
