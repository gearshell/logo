package logo

import (
	"image"
	"image/color"
)

func fillTransparentBlack(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			if a < 128 {
				dst.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				dst.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
			}
		}
	}
	return dst
}
