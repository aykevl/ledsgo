# Color utilities for LED animations

This package is a collection of some utility functions for working with
color. It is primarily intended for LED animations on microcontrollers using
[TinyGo](https://tinygo.org/), therefore it has been optimized for devices
without FPU.

It is inspired by [FastLED](http://fastled.io/) but does not implement any
drivers for LED strips to keep development focused on fast animations.

## Noise functions

This package contains a number of Simplex noise functions.
[Simplex noise](https://en.wikipedia.org/wiki/Simplex_noise) is very similar
to Perlin noise and produces naturally looking gradients as you might
encounter in nature. It is commonly used as a building block for animations,
especially in procedurally generated games.

Be warned that Simplex noise is
[patented](https://patents.google.com/patent/US6867776) (set to expire on
2022-01-18) so use at your own risk for computer graphics. This patent may or
may not apply to LED animations, I don't know.

## License

This package is licensed under a BSD-style license, see the LICENSE file for
details.
