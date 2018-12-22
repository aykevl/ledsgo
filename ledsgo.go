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
