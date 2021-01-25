package mandelbrot

import (
	"image"
	"mandelbrot/fractal"
	"mandelbrot/palette"
	"math/big"
	"math/cmplx"
	"sync"
	"sync/atomic"
)

// MandelbrotFloat64 is a mandelbrot fractal generator that uses float64 numbers for calculations
type Float64 struct {
	iterations int
	threshold  float32
}

func NewFloat64Default() *Float64 {
	return NewFloat64(DefaultIterations, DefaultThreshold)
}

func NewFloat64(iterations int, threshold float32) *Float64 {
	ret := &Float64{
		iterations: iterations,
		threshold:  threshold,
	}

	return ret
}

// Generation function
func (f *Float64) Generate(
	target *image.RGBA,
	cx, cy, scale *big.Float,
	physicalWidth, physicalHeight *big.Float,
	reportingFunc fractal.ProgressReportingFunc,
	doneFunc fractal.DoneFunc,
) {
	go func() {
		x, _ := cx.Float64()
		y, _ := cy.Float64()
		scalef64, _ := scale.Float64()

		// Calculate physical width and height
		physWidthF64, _ := physicalWidth.Float64()
		physHeightF64, _ := physicalHeight.Float64()

		physWidth := physWidthF64 * scalef64
		physHeight := physHeightF64 * scalef64

		width := target.Rect.Max.X
		height := target.Rect.Max.Y

		// Scale physical bounds
		physMinX := x - (physWidth / 2)
		physMinY := y - (physHeight / 2)

		// Calculate pixel-to-physical scale
		scaleX := physWidth / float64(width)
		scaleY := physHeight / float64(height)

		pal := palette.CreatePaletteGrayscaleRecursive(256)

		wg := sync.WaitGroup{}

		// Counter will hold number of completed lines. Need for progress reporting
		var linesDone int32
		linesDonePtr := &linesDone

		// (x, y) - are pixel coords
		for y := 0; y < height; y++ {
			// (physX, physY) - are physical coordinates
			physY := float64(y)*scaleY + physMinY
			wg.Add(1)
			go func(y int, physY float64) {
				for x := 0; x < width; x++ {
					physX := float64(x)*scaleX + physMinX

					// get fractal value at the point
					value := mandelbrotComplex128(complex(physX, physY), f.iterations, f.threshold)

					// convert it to the color and set pixel color
					target.Set(x, y, pal[int(float32(len(pal))*value)])
				}
				atomic.AddInt32(linesDonePtr, 1)
				reportingFunc(float32(atomic.LoadInt32(linesDonePtr)) / float32(height))
				wg.Done()
			}(y, physY)
		}

		wg.Wait()
		doneFunc()
	}()
}

// Calculate mandelbrot set value in given point
func mandelbrotComplex128(c complex128, iterations int, threshold float32) float32 {
	ret := complex(0, 0)

	for i := 0; i < iterations; i++ {
		ret = ret*ret + c
		if float32(cmplx.Abs(ret)) > threshold {
			return float32(i) / float32(iterations)
		}
	}

	return 0
}
