package main

import (
	"fmt"
	"image"
	"image/png"
	big2 "math/big"
	"os"
	"sync"

	"mandelbrot/fractal"
	"mandelbrot/palette"
	"mandelbrot/big"
)

const precision = 16

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("usage: mbrot <cx> <cy> <scale>\n")
		return
	}

	cx := big.MustParseBigFloat(args[1], precision)
	cy := big.MustParseBigFloat(args[2], precision)
	scale := big.MustParseBigFloat(args[3], precision)

	// pixel bounds
	screenMinX := 0
	screenMaxX := 640
	screenMinY := 0
	screenMaxY := 480

	screenWidth := screenMaxX - screenMinX
	screenHeight := screenMaxY - screenMinY

	initWidth := big2.NewFloat(3)
	initHeight := big2.NewFloat(2)

	scaledWidth := big2.NewFloat(0).Mul(initWidth, scale)
	scaledHeight := big2.NewFloat(0).Mul(initHeight, scale)

	half := big2.NewFloat(0.5)

	// physical bounds
	scaledWidthHalf := big2.NewFloat(0).Copy(scaledWidth)
	scaledWidthHalf = scaledWidthHalf.Mul(scaledWidthHalf, half)

	physMinX := big2.NewFloat(0).Copy(cx)
	physMinX.Sub(physMinX, scaledWidthHalf)

	scaledHeightHalf := big2.NewFloat(0).Copy(scaledHeight)
	scaledHeightHalf = scaledHeightHalf.Mul(scaledHeightHalf, half)

	physMinY := big2.NewFloat(0).Copy(cy)
	physMinY.Sub(physMinY, scaledHeightHalf)

	// pixel-to-physical scale
	scaleX := big2.NewFloat(0).Quo(scaledWidth, big2.NewFloat(float64(screenWidth)))
	scaleY := big2.NewFloat(0).Quo(scaledHeight, big2.NewFloat(float64(screenHeight)))

	pal := palette.CreatePaletteGrayscaleRecursive(256)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: screenMinX, Y: screenMinY},
		Max: image.Point{X: screenMaxX, Y: screenMaxY},
	})

	var wg sync.WaitGroup

	for y := screenMinY; y < screenMaxY; y++ {
		// (physX, physY) - are physical coordinates
		physY := big2.NewFloat(float64(y))
		physY.Mul(physY, scaleY)
		physY.Add(physY, physMinY)

		wg.Add(1)
		go func(y int) {
			for x := screenMinX; x < screenMaxX; x++ {
				physX := big2.NewFloat(float64(x))
				physX.Mul(physX, scaleX)
				physX.Add(physX, physMinX)

				// get fractal value at the point
				value := fractal.MandelbrotBig(physX, physY)

				// convert it to the color and set pixel color
				img.Set(x, screenMaxY - y - 1, pal[int(float64(len(pal)) * value)])
			}
			wg.Done()
		}(y)
	}

	wg.Wait()

	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}

