package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"

	"mandelbrot/fractal"
	"mandelbrot/palette"
)

func main_simple() {
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("usage: mbrot <cx> <cy> <scale>\n")
		return
	}

	cx, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		panic(err)
	}

	cy, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		panic(err)
	}

	scale, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		panic(err)
	}

	// pixel bounds
	screenMinX := 0
	screenMaxX := 640
	screenMinY := 0
	screenMaxY := 480

	screenWidth := screenMaxX - screenMinX
	screenHeight := screenMaxY - screenMinY

	initWidth := 3.0
	initHeight := 2.0

	scaledWidth := initWidth * scale
	scaledHeight := initHeight * scale

	// physical bounds
	physMinX := cx - (scaledWidth / 2)
	physMinY := cy - (scaledHeight / 2)

	// pixel-to-physical scale
	scaleX := scaledWidth / float64(screenWidth)
	scaleY := scaledHeight / float64(screenHeight)

	pal := palette.CreatePaletteGrayscaleRecursive(256)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: screenMinX, Y: screenMinY},
		Max: image.Point{X: screenMaxX, Y: screenMaxY},
	})

	_, _ = fmt.Fprintf(os.Stderr, "scaled width : %.10f\n", scaledWidth)
	_, _ = fmt.Fprintf(os.Stderr, "scaled height: %.10f\n", scaledHeight)
	_, _ = fmt.Fprintf(os.Stderr, "scale X      : %.10f\n", scaleX)
	_, _ = fmt.Fprintf(os.Stderr, "scale Y      : %.10f\n", scaleY)

	// (x, y) - are pixel coords
	for y := screenMinY; y < screenMaxY; y++ {
		// (physX, physY) - are physical coordinates
		physY := float64(y)*scaleY + physMinY
		for x := screenMinX; x < screenMaxX; x++ {
			physX := float64(x)*scaleX + physMinX

			// get fractal value at the point
			value := fractal.Mandelbrot(complex(physX, physY))

			// convert it to the color and set pixel color
			img.Set(x, screenMaxY - y - 1, pal[int(float64(len(pal)) * value)])
		}
	}

	err = png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}

