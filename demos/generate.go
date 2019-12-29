// +build none

// This file is used in `go generate` to update the images directory.

package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/aykevl/ledsgo/demos"
	"github.com/kettek/apng"
)

const (
	width  = 32
	height = 32
	scale  = 4
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "provide exactly one argumet: the directory to store the resulting images")
		os.Exit(1)
	}
	dir := os.Args[1] // destination directory

	saveAnimation(demos.Fire, filepath.Join(dir, "fire.png"))
}

func saveAnimation(draw func(demos.Displayer, time.Time), path string) {
	fmt.Println("generating:", path)
	const frameRate = 30 // frames per second
	const duration = 3   // seconds
	const frames = frameRate * duration
	const frameTime = time.Second / frameRate
	display := &imageDisplayer{
		Scale: scale,
	}
	for i := 0; i < frames; i++ {
		t := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Add(frameTime * time.Duration(i))
		display.NextFrame(frameRate)
		draw(display, t)
	}

	// Write image to directory.
	err := display.Save(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to save image:", err)
	}
}

type imageDisplayer struct {
	frame *image.RGBA
	img   apng.APNG
	Scale int
}

func (d *imageDisplayer) Size() (width, height int16) {
	rect := d.frame.Bounds()
	return int16(rect.Max.X / d.Scale), int16(rect.Max.Y / d.Scale)
}

func (d *imageDisplayer) SetPixel(x, y int16, c color.RGBA) {
	c = gammaEncode(c)
	for ix := int(x) * d.Scale; ix < int(x+1)*d.Scale; ix++ {
		for iy := int(y) * d.Scale; iy < int(y+1)*d.Scale; iy++ {
			d.frame.Set(ix, int(iy), c)
		}
	}
}

func (d *imageDisplayer) Display() error {
	return nil // nop
}

func (d *imageDisplayer) NextFrame(frameRate uint16) {
	if d.frame != nil {
		d.img.Frames = append(d.img.Frames, apng.Frame{
			Image:            d.frame,
			DelayNumerator:   1,
			DelayDenominator: frameRate,
		})
	}
	d.frame = image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))
}

func (d *imageDisplayer) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return apng.Encode(f, d.img)
}

// gammaEncode encodes a color value into sRGB, the color space often used in
// computer graphics.
func gammaEncode(c color.RGBA) color.RGBA {
	const gamma = 1 / 2.2
	return color.RGBA{
		R: uint8(math.Pow(float64(c.R)/255, gamma) * 255.5),
		G: uint8(math.Pow(float64(c.G)/255, gamma) * 255.5),
		B: uint8(math.Pow(float64(c.B)/255, gamma) * 255.5),
		A: c.A,
	}
}
