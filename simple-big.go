package main

import (
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	big2 "math/big"
	"os"
	"sync"

	"mandelbrot/big"
	"mandelbrot/fractal"
	"mandelbrot/palette"
)

const precision = 1024

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

	initWidth := big2.NewFloat(3).SetPrec(precision)
	initHeight := big2.NewFloat(2).SetPrec(precision)

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
		// physX := y * scaleY + physMinY
		physY := big2.NewFloat(float64(y))
		physY.Mul(physY, scaleY)
		physY.Add(physY, physMinY)

		wg.Add(1)
		go func(y int, physY *big2.Float) {
			physX := big2.NewFloat(0).Set(physMinX)
			for x := screenMinX; x < screenMaxX; x++ {
				// get fractal value at the point
				value := fractal.MandelbrotBig(physX, physY)

				// convert it to the color and set pixel color
				img.Set(x, screenMaxY - y - 1, pal[int(float64(len(pal)) * value)])

				physX.Add(physX, scaleX)
			}
			wg.Done()
		}(y, physY)
	}

	wg.Wait()

	lineHeight := 13
	strX := 2
	strY := 0

	strY+=lineHeight
	drawString(img, strX, strY, fmt.Sprintf("cx    : %+.30f", cx))

	strY+=lineHeight
	drawString(img, strX, strY, fmt.Sprintf("cy    : %+.30f", cy))

	strY+=lineHeight
	drawString(img, strX, strY, fmt.Sprintf("xmin  : %+.30f", physMinX))

	strY+=lineHeight
	drawString(img, strX, strY, fmt.Sprintf("width : %+.30f", scaledWidth))

	strY+=lineHeight
	drawString(img, strX, strY, fmt.Sprintf("ymin  : %+.30f", physMinY))

	strY+=lineHeight
	drawString(img, strX, strY, fmt.Sprintf("height: %+.30f", scaledHeight))

	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}

func drawString(img draw.Image, x int, y int, str string) {
	dot := fixed.P(x, y)
	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{B: 0xFF, A: 0xFF}),
		Face: basicfont.Face7x13,
		Dot:  dot,
	}
	d.DrawString(str)

}
