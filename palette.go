package ledsgo

import (
	"image/color"
)

// Palette16 is a 16-color palette on a continuous scale, from which a color can
// be picked.
type Palette16 [16]color.RGBA

// ColorAt returns a color from the palette at the 16-bit index (0..65535)
// position. Colors not exactly from one position are interpolated. Colors close
// to the top wrap around to the bottom, so you can imagine the palette as being
// a custom color ring akin to a hue color ring in a HSV color picker.
//
// This method is similar to ColorFromPalette in FastLED except that it uses a
// 16-bit index instead of a 8-bit index for better color accuracy.
func (p *Palette16) ColorAt(position uint16) color.RGBA {
	index := position >> 12
	blendPosition := uint8(position >> 4)

	bottom := p[index]

	var top color.RGBA
	if index >= 15 {
		top = p[0]
	} else {
		top = p[index+1]
	}

	return color.RGBA{
		R: scale8(bottom.R, 255-blendPosition) + scale8(top.R, blendPosition),
		G: scale8(bottom.G, 255-blendPosition) + scale8(top.G, blendPosition),
		B: scale8(bottom.B, 255-blendPosition) + scale8(top.B, blendPosition),
		A: 0xff,
	}
}

// A number of palettes copied from FastLED, see:
// https://github.com/FastLED/FastLED/blob/master/colorpalettes.cpp
var (
	CloudColors = Palette16{
		Blue,
		DarkBlue,
		DarkBlue,
		DarkBlue,

		DarkBlue,
		DarkBlue,
		DarkBlue,
		DarkBlue,

		Blue,
		DarkBlue,
		SkyBlue,
		SkyBlue,

		LightBlue,
		White,
		LightBlue,
		SkyBlue,
	}

	LavaColors = Palette16{
		Black,
		Maroon,
		Black,
		Maroon,

		DarkRed,
		Maroon,
		DarkRed,

		DarkRed,
		DarkRed,
		Red,
		Orange,

		White,
		Orange,
		Red,
		DarkRed,
	}

	OceanColors = Palette16{
		MidnightBlue,
		DarkBlue,
		MidnightBlue,
		Navy,

		DarkBlue,
		MediumBlue,
		SeaGreen,
		Teal,

		CadetBlue,
		Blue,
		DarkCyan,
		CornflowerBlue,

		Aquamarine,
		SeaGreen,
		Aqua,
		LightSkyBlue,
	}

	ForestColors = Palette16{
		DarkGreen,
		DarkGreen,
		DarkOliveGreen,
		DarkGreen,

		Green,
		ForestGreen,
		OliveDrab,
		Green,

		SeaGreen,
		MediumAquamarine,
		LimeGreen,
		YellowGreen,

		LightGreen,
		LawnGreen,
		MediumAquamarine,
		ForestGreen,
	}

	// HSV Rainbow
	RainbowColors = Palette16{
		color.RGBA{0xFF, 0x00, 0x00, 0xFF}, color.RGBA{0xD5, 0x2A, 0x00, 0xFF}, color.RGBA{0xAB, 0x55, 0x00, 0xFF}, color.RGBA{0xAB, 0x7F, 0x00, 0xFF},
		color.RGBA{0xAB, 0xAB, 0x00, 0xFF}, color.RGBA{0x56, 0xD5, 0x00, 0xFF}, color.RGBA{0x00, 0xFF, 0x00, 0xFF}, color.RGBA{0x00, 0xD5, 0x2A, 0xFF},
		color.RGBA{0x00, 0xAB, 0x55, 0xFF}, color.RGBA{0x00, 0x56, 0xAA, 0xFF}, color.RGBA{0x00, 0x00, 0xFF, 0xFF}, color.RGBA{0x2A, 0x00, 0xD5, 0xFF},
		color.RGBA{0x55, 0x00, 0xAB, 0xFF}, color.RGBA{0x7F, 0x00, 0x81, 0xFF}, color.RGBA{0xAB, 0x00, 0x55, 0xFF}, color.RGBA{0xD5, 0x00, 0x2B, 0xFF},
	}

	// HSV Rainbow colors with alternatating stripes of black
	RainbowStripeColors = Palette16{
		color.RGBA{0xFF, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0xAB, 0x55, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
		color.RGBA{0xAB, 0xAB, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0xFF, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
		color.RGBA{0x00, 0xAB, 0x55, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0xFF, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
		color.RGBA{0x55, 0x00, 0xAB, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF}, color.RGBA{0xAB, 0x00, 0x55, 0xFF}, color.RGBA{0x00, 0x00, 0x00, 0xFF},
	}

	// HSV color ramp: blue purple ping red orange yellow (and back)
	// Basically, everything but the greens, which tend to make
	// people's skin look unhealthy.  This palette is good for
	// lighting at a club or party, where it'll be shining on people.
	PartyColors = Palette16{
		color.RGBA{0x55, 0x00, 0xAB, 0xFF}, color.RGBA{0x84, 0x00, 0x7C, 0xFF}, color.RGBA{0xB5, 0x00, 0x4B, 0xFF}, color.RGBA{0xE5, 0x00, 0x1B, 0xFF},
		color.RGBA{0xE8, 0x17, 0x00, 0xFF}, color.RGBA{0xB8, 0x47, 0x00, 0xFF}, color.RGBA{0xAB, 0x77, 0x00, 0xFF}, color.RGBA{0xAB, 0xAB, 0x00, 0xFF},
		color.RGBA{0xAB, 0x55, 0x00, 0xFF}, color.RGBA{0xDD, 0x22, 0x00, 0xFF}, color.RGBA{0xF2, 0x00, 0x0E, 0xFF}, color.RGBA{0xC2, 0x00, 0x3E, 0xFF},
		color.RGBA{0x8F, 0x00, 0x71, 0xFF}, color.RGBA{0x5F, 0x00, 0xA1, 0xFF}, color.RGBA{0x2F, 0x00, 0xD0, 0xFF}, color.RGBA{0x00, 0x07, 0xF9, 0xFF},
	}

	// Approximate "black body radiation" palette, akin to
	// the FastLED 'HeatColor' function.
	// Recommend that you use values 0-240 rather than
	// the usual 0-255, as the last 15 colors will be
	// 'wrapping around' from the hot end to the cold end,
	// which looks wrong.
	HeatColors = Palette16{
		color.RGBA{0x00, 0x00, 0x00, 0xFF},
		color.RGBA{0x33, 0x00, 0x00, 0xFF}, color.RGBA{0x66, 0x00, 0x00, 0xFF}, color.RGBA{0x99, 0x00, 0x00, 0xFF}, color.RGBA{0xCC, 0x00, 0x00, 0xFF}, color.RGBA{0xFF, 0x00, 0x00, 0xFF},
		color.RGBA{0xFF, 0x33, 0x00, 0xFF}, color.RGBA{0xFF, 0x66, 0x00, 0xFF}, color.RGBA{0xFF, 0x99, 0x00, 0xFF}, color.RGBA{0xFF, 0xCC, 0x00, 0xFF}, color.RGBA{0xFF, 0xFF, 0x00, 0xFF},
		color.RGBA{0xFF, 0xFF, 0x33, 0xFF}, color.RGBA{0xFF, 0xFF, 0x66, 0xFF}, color.RGBA{0xFF, 0xFF, 0x99, 0xFF}, color.RGBA{0xFF, 0xFF, 0xCC, 0xFF}, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	}
)
