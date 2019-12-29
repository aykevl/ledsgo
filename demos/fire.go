package demos

import (
	"image/color"
	"time"

	"github.com/aykevl/ledsgo"
)

// Fire shows an animation that looks somewhat like fire. The 'now' time
// indicates which instance of the animation is generated.
func Fire(display Displayer, now time.Time) {
	width, height := display.Size()
	const pointsPerCircle = 12 // how many LEDs there are per turn of the torch
	const speed = 12           // higher means faster
	var cooling = 256 / height // higher means faster cooling
	var detail = 12800 / width // higher means more detailed flames
	for x := int16(0); x < width; x++ {
		for y := int16(0); y < width; y++ {
			heat := ledsgo.Noise2(int32(y*detail)+int32((now.UnixNano()>>20)*speed), int32(x*detail))/256 + 128
			heat -= int16((height-1)-y) * cooling
			if heat < 0 {
				heat = 0
			}
			display.SetPixel(x, y, heatMap(uint8(heat)))
		}
	}
}

func heatMap(index uint8) color.RGBA {
	if index < 128 {
		return color.RGBA{index * 2, 0, 0, 255}
	}
	if index < 224 {
		return color.RGBA{255, uint8(uint32(index-128) * 8 / 3), 0, 255}
	}
	return color.RGBA{255, 255, (index - 224) * 8, 255}
}
