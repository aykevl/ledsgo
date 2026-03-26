package ledsgo

// This file implements the "improved Perlin" algorithms, copied from FastLED.
// They are a bit faster than the Simplex noise algorithms.
// See:
// https://github.com/FastLED/FastLED/blob/master/src/noise.cpp.hpp

var fastled_P = [...]uint8{
	151, 160, 137, 91, 90, 15, 131, 13, 201, 95, 96, 53, 194, 233, 7, 225,
	140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23, 190, 6, 148,
	247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32,
	57, 177, 33, 88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175,
	74, 165, 71, 134, 139, 48, 27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
	60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244, 102, 143, 54,
	65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169,
	200, 196, 135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64,
	52, 217, 226, 250, 124, 123, 5, 202, 38, 147, 118, 126, 255, 82, 85, 212,
	207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
	119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104,
	218, 246, 97, 228, 251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241,
	81, 51, 145, 235, 249, 14, 239, 107, 49, 192, 214, 31, 181, 199, 106, 157,
	184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
	222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
	151}

// Scale a 16-bit unsigned value by an 16-bit value.
func scale16(i uint16, scale uint16) uint16 {
	return uint16((uint32(i) * (1 + uint32(scale))) / 65536)
}

func ease16InOutQuad(i uint16) uint16 {
	// This is the legacy version, there is a slightly more accurate version in fl/ease.cpp
	// with fl::easeInOutQuad16. However the difference is minimal.
	//
	// 16-bit quadratic ease-in / ease-out function
	j := i
	if j&0x8000 != 0 {
		j = 65535 - j
	}
	jj := scale16(j, j)
	jj2 := jj << 1
	if i&0x8000 != 0 {
		jj2 = 65535 - jj2
	}
	return jj2
}

func avg15(i, j int16) int16 {
	return (i >> 1) + (j >> 1) + (i & 0x1)
}

func grad16(hash uint8, x, y int16) int16 {
	hash = hash & 7
	var u, v int16
	if hash < 4 {
		u = x
		v = y
	} else {
		u = y
		v = x
	}
	if hash&1 != 0 {
		u = -u
	}
	if hash&2 != 0 {
		v = -v
	}

	return avg15(u, v)
}

// Linear interpolation between two signed 15-bit values, with 8-bit fraction.
func lerp15by16(a, b int16, frac uint16) int16 {
	var result int16
	if b > a {
		delta := uint16(b - a)
		scaled := scale16(delta, frac)
		result = a + int16(scaled)
	} else {
		delta := uint16(a - b)
		var scaled uint16 = scale16(delta, frac)
		result = a - int16(scaled)
	}
	return result
}

func inoise16_raw(x, y uint32) int16 {
	// Find the unit cube containing the point
	X := uint8(x >> 16)
	Y := uint8(y >> 16)

	// Hash cube corner coordinates
	A := fastled_P[X] + Y
	AA := fastled_P[A]
	AB := fastled_P[uint32(A)+1]
	B := fastled_P[uint32(X)+1] + Y
	BA := fastled_P[B]
	BB := fastled_P[uint32(B)+1]

	// Get the relative position of the point in the cube
	u := uint16(x)
	v := uint16(y)

	// Get a signed version of the above for the grad function
	xx := int16((u >> 1) & 0x7FFF)
	yy := int16((v >> 1) & 0x7FFF)
	N := int16(-0x8000)

	u = ease16InOutQuad(u)
	v = ease16InOutQuad(v)

	X1 := lerp15by16(grad16(fastled_P[AA], xx, yy), grad16(fastled_P[BA], xx-N, yy), u)
	X2 := lerp15by16(grad16(fastled_P[AB], xx, yy-N), grad16(fastled_P[BB], xx-N, yy-N), u)

	return lerp15by16(X1, X2, v)
}

// 2D Perlin noise, copied from FastLED (inoise16).
func Perlin2(x, y uint32) uint16 {
	ans := int32(inoise16_raw(x, y))
	ans = ans + 17308
	pan := uint32(ans)
	// pan = (ans * 242L) >> 7.  That's the same as:
	// pan = (ans * 484L) >> 8.  And this way avoids a 7X four-byte shift-loop on AVR.
	// Identical math, except for the highest bit, which we don't care about anyway,
	// since we're returning the 'middle' 16 out of a 32-bit value anyway.
	pan *= 484
	return uint16(pan >> 8)
}
