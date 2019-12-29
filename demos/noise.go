package demos

import (
	"time"

	"github.com/aykevl/ledsgo"
)

// Noise shows noise mapped to a rainbow function. The 'now' time indicates
// which instance of the animation is generated.
func Noise(display Displayer, now time.Time) {
	const spread = 6 // higher means the noise gets more detailed
	const speed = 20 // higher means slower
	width, height := display.Size()
	for x := int16(0); x < width; x++ {
		for y := int16(0); y < height; y++ {
			hue := uint16(ledsgo.Noise3(int32(now.UnixNano()>>speed), int32(x<<spread), int32(y<<spread))) * 2
			display.SetPixel(x, y, ledsgo.Color{hue, 0xff, 0xff}.Spectrum())
		}
	}
}
