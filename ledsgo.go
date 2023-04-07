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
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
}

// Rainbow color conversion. This function was copied from the FastLED library,
// for details see
// https://github.com/FastLED/FastLED/wiki/FastLED-HSV-Colors#color-map-rainbow-vs-spectrum.
// The main difference is that yellow is somewhat brighter, which may look
// better in LEDs.
func (c Color) Rainbow() color.RGBA {
	// Yellow has a higher inherent brightness than any other color; 'pure'
	// yellow is perceived to be 93% as bright as white.  In order to make
	// yellow appear the correct relative brightness, it has to be rendered
	// brighter than all other colors.
	// Level Y1 is a moderate boost, the default.
	// Level Y2 is a strong boost.
	const Y1 = true
	const Y2 = false

	hue := uint8(c.H >> 8)
	sat := c.S
	val := c.V

	offset := hue & 0x1F // 0..31

	offset8 := offset * 8

	third := scale8(offset8, (256 / 3)) // max = 85

	var r, g, b uint8

	if hue&0x80 == 0 {
		// 0XX
		if hue&0x40 == 0 {
			// 00X
			//section 0-1
			if hue&0x20 == 0 {
				// 000
				//case 0: // R -> O
				r = 255 - third
				g = third
				b = 0
			} else {
				// 001
				//case 1: // O -> Y
				if Y1 {
					r = 171
					g = 85 + third
					b = 0
				}
				if Y2 {
					r = 170 + third
					twothirds := scale8(offset8, ((256 * 2) / 3)) // max=170
					g = 85 + twothirds
					b = 0
				}
			}
		} else {
			//01X
			// section 2-3
			if hue&0x20 == 0 {
				// 010
				//case 2: // Y -> G
				if Y1 {
					twothirds := scale8(offset8, ((256 * 2) / 3)) // max=170
					r = 171 - twothirds
					g = 170 + third
					b = 0
				}
				if Y2 {
					r = 255 - offset8
					g = 255
					b = 0
				}
			} else {
				// 011
				// case 3: // G -> A
				r = 0
				g = 255 - third
				b = third
			}
		}
	} else {
		// section 4-7
		// 1XX
		if hue&0x40 == 0 {
			// 10X
			if hue&0x20 == 0 {
				// 100
				//case 4: // A -> B
				r = 0
				//uint8_t twothirds = (third << 1);
				twothirds := scale8(offset8, ((256 * 2) / 3)) // max=170
				g = 171 - twothirds                           //170?
				b = 85 + twothirds
			} else {
				// 101
				//case 5: // B -> P
				r = third
				g = 0
				b = 255 - third
			}
		} else {
			if hue&0x20 == 0 {
				// 110
				//case 6: // P -- K
				r = 85 + third
				g = 0
				b = 171 - third
			} else {
				// 111
				//case 7: // K -> R
				r = 170 + third
				g = 0
				b = 85 - third
			}
		}
	}

	// Scale down colors if we're desaturated at all
	// and add the brightness_floor to r, g, and b.
	if sat != 255 {
		if sat == 0 {
			r = 255
			b = 255
			g = 255
		} else {
			//nscale8x3_video( r, g, b, sat);
			if r != 0 {
				r = scale8(r, sat)
			}
			if g != 0 {
				g = scale8(g, sat)
			}
			if b != 0 {
				b = scale8(b, sat)
			}

			desat := 255 - sat
			desat = scale8(desat, desat)

			brightness_floor := desat
			r += brightness_floor
			g += brightness_floor
			b += brightness_floor
		}
	}

	// Now scale everything down if we're at value < 255.
	if val != 255 {
		val = scale8_video(val, val)
		if val == 0 {
			r = 0
			g = 0
			b = 0
		} else {
			// nscale8x3_video( r, g, b, val);
			if r != 0 {
				r = scale8(r, val)
			}
			if g != 0 {
				g = scale8(g, val)
			}
			if b != 0 {
				b = scale8(b, val)
			}
		}
	}

	return color.RGBA{r, g, b, 0xff}
}

// Scale one byte by a second one, which is treated as the numerator of a
// fraction whose denominator is 256.
// In other words, it computes i * (scale / 256)
func scale8(i, scale uint8) uint8 {
	return uint8((uint16(i) * (1 + uint16(scale))) >> 8)
}

// The "video" version of scale8 guarantees that the output will only be zero if
// one or both of the inputs are zero. If both inputs are non-zero, the output
// is guaranteed to be non-zero. This makes for better 'video'/LED dimming, at
// the cost of several additional cycles.
func scale8_video(i, scale uint8) uint8 {
	result := ((int(i) * int(scale)) >> 8)
	if i != 0 && scale != 0 {
		result += 1
	}
	return uint8(result)
}

// Blend the top value into the bottom value, with the given alpha value.
func blend(bottom, top, topAlpha uint8) uint8 {
	return uint8((int(bottom)*(256-int(topAlpha)) + int(top)*int(topAlpha)) >> 8)
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
