// Package patterns implements a number of animations for small LED displays.
package demos

//go:generate go run generate.go ./images

import "image/color"

// Displayer is a surface that can be drawn upon.
type Displayer interface {
	// Size returns the current size of the display.
	Size() (x, y int16)

	// SetPizel modifies the internal buffer.
	SetPixel(x, y int16, c color.RGBA)
}
