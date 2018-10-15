package ledsgo

// The LEDs of a LED strip. Colors are represented in RGB form: 0x00rrggbb.
type Strip []uint32

// Fill the LED strip with a color range, using the HSV spectrum conversion.
// The start color is the color for the first LED. All other colors have the
// same saturation and value but increased (and wrapped) hue.
func (s Strip) FillSpectrum(start Color, hueinc uint16) {
	for i := range s {
		s[i] = start.Spectrum()
		start.H += hueinc
	}
}

// FillSolid sets all colors to the given value.
func (s Strip) FillSolid(color uint32) {
	for i := range s {
		s[i] = color
	}
}
