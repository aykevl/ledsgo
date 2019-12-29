package ledsgo

import (
	"image/color"
)

// Color encodes a HSV color.
//
// The hue is 16-bits to get better looking colors, as HSV→RGB conversions
// generally can use more than 8 hue bits for their conversion. Saturation and
// value are both just 8 bits because saturation is not that often used and
// value does not gain much precision with extra bits. This encoding has been
// chosen to have the best colors while still fitting in 32 bits.
type Color struct {
	H uint16 // hue
	S uint8  // saturation
	V uint8  // value
}

// Return the RGB version of this color, calculated with the common HSV
// conversion:
// https://github.com/FastLED/FastLED/wiki/FastLED-HSV-Colors
func (c Color) Spectrum() color.RGBA {
	sectionWidth := uint32((1<<16)/3 + 1) // one third of the hue space
	section := uint32(c.H) / sectionWidth
	colorValue := (uint32(c.H) - section*sectionWidth) * 256 / sectionWidth
	var r, g, b uint32
	switch section {
	case 0:
		r = 0xff - colorValue
		g = colorValue
	case 1:
		g = 0xff - colorValue
		b = colorValue
	case 2:
		b = 0xff - colorValue
		r = colorValue
	}
	sat := (0xff - uint32(c.S)) * uint32(c.V) / 256 / 3
	r = r*uint32(c.V)*(uint32(c.S))/(1<<16) + sat
	g = g*uint32(c.V)*(uint32(c.S))/(1<<16) + sat
	b = b*uint32(c.V)*(uint32(c.S))/(1<<16) + sat
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0}
}

// ApplyAlpha scales the color with the given alpha. It can be used to reduce
// the intensity of a given color. The color is assumed to be linear, not sRGB.
func ApplyAlpha(c color.RGBA, alpha uint8) color.RGBA {
	return color.RGBA{
		R: uint8(uint32(c.R) * uint32(alpha) / 0xff),
		G: uint8(uint32(c.G) * uint32(alpha) / 0xff),
		B: uint8(uint32(c.B) * uint32(alpha) / 0xff),
		A: uint8(uint32(c.A) * uint32(alpha) / 0xff),
	}
}

// Blend blends two colors together, assuming the colors are linear (not sRGB).
// The bottom alpha is assumed to be 0xff. The top alpha is used to blend the
// two colors together.
func Blend(bottom, top color.RGBA) color.RGBA {
	return color.RGBA{
		R: uint8((uint32(bottom.R)*uint32(255-top.A))/255 + uint32(top.R)),
		G: uint8((uint32(bottom.G)*uint32(255-top.A))/255 + uint32(top.G)),
		B: uint8((uint32(bottom.B)*uint32(255-top.A))/255 + uint32(top.B)),
		A: 255,
	}
}

// Sqrt is a fast integer square root function that returns a result at most one
// off of the intended result. It can be used in place of a floating-point
// square root function. Negative values are accepted and result in 0.
func Sqrt(x int) int {
	// Original code copied from:
	// https://stackoverflow.com/questions/34187171/fast-integer-square-root-approximation/#34187992

	if x < 0 {
		// Avoid division by 0.
		// Also avoid negative numbers. Returning 0 is incorrect but should be
		// fine for most animation purposes.
		return 0
	}

	a := 1024
	b := x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2
	b = x / a
	a = (a + b) / 2

	return a
}
