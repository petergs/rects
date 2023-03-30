package main

import (
	"image"
	"image/color"
)

func drawRect(dest *image.RGBA, w, h, cx, cy int, clr color.NRGBA) {
	//fmt.Printf("cx: %d, cy: %d\n", cx, cy)
	cx = cx - w/2
	cy = cy - h/2
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dest.Set(cx+x, cy+y, clr)
		}
	}
}

func drawBackground(dest *image.RGBA, clr color.NRGBA) {
	max := dest.Bounds().Max

	for x := 0; x < max.X; x++ {
		for y := 0; y < max.Y; y++ {
			dest.Set(x, y, clr)
		}
	}
}

func drawManyRects(s *settings, dest *image.RGBA) {

	// draw positions
	var cx int
	var cy int

	// jitter offsets
	var wj int
	var hj int

	// iterate over all pixels
	// check if there should be a rect centered there
	// draw. them. rects.
	for x := 1; x < s.canvasWidth; x++ {
		for y := 1; y < s.canvasHeight; y++ {
			if y%s.distY == 0 && x%s.distX == 0 {
				cx = x + randRange(-1*s.jitterX, s.jitterX)
				cy = y + randRange(-1*s.jitterY, s.jitterY)

				wj = randRange(-1*s.jitterWidth, s.jitterWidth)
				if s.preserveSquare {
					hj = wj
				} else {
					hj = randRange(-1*s.jitterHeight, s.jitterHeight)
				}

				clr := randColor(s.rectColors)
				if s.fromSrcImage {
					r, g, b, a := dest.At(cx, cy).RGBA()
					colorAt := color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
					if colorAt != clr {
						drawRect(dest, s.rectWidth+wj, s.rectHeight+hj, cx, cy, clr)
					}
				} else {
					drawRect(dest, s.rectWidth+wj, s.rectHeight+hj, cx, cy, clr)
				}

			}
		}
	}

}
