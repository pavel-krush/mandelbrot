package palette

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func CreatePaletteGrayscaleLinear(values int) color.Palette {
	ret := make([]color.Color, values)
	step := float64(256 / values)

	for i := 0; i < values; i++ {
		val := uint8(float64(i) * step)
		ret[i] = color.RGBA{
			R: val,
			G: val,
			B: val,
			A: 255,
		}
	}

	return ret
}

func CreatePaletteGrayscaleRecursive(values int) color.Palette {
	ret := make([]color.Color, values)
	_createPaletteGrayscaleRecursive(ret, 0, values-1)
	return ret
}

func _createPaletteGrayscaleRecursive(palette color.Palette, l int, r int) {
	//_,_ = fmt.Fprintf(os.Stderr, "l: %d, r: %d\n", l, r)

	if l == r {
		palette[l] = color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		}
		return
	}

	mid := l + (r - l) / 2
	//_,_ = fmt.Fprintf(os.Stderr, "mid: %d\n", mid)

	step := 255 / float64(mid-l)
	//_,_ = fmt.Fprintf(os.Stderr, "step: %f\n", step)

	for i := l; i <= mid; i++ {
		val := uint8(float64(i - l) * step)
		palette[i] = color.RGBA{
			R: val,
			G: val,
			B: val,
			A: 255,
		}
	}

	_createPaletteGrayscaleRecursive(palette, mid+1, r)
}

func DrawPalette(palette color.Palette) {
	width := 640
	height := 32

	colorsCount := len(palette)
	_,_ = fmt.Fprintf(os.Stderr, "colors count: %d\n", colorsCount)

	colorWidth := float64(width) / float64(colorsCount)
	_,_ = fmt.Fprintf(os.Stderr, "color cell width: %f\n", colorWidth)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorValue := float64(x) / colorWidth
			colorNum := int(colorValue)
			img.Set(x, y, palette[colorNum])
		}
	}

	png.Encode(os.Stdout, img)
}

func DumpPalette(palette color.Palette) {
	for i := range palette {
		fmt.Printf("%d %+v\n", i, palette[i])
	}
}
