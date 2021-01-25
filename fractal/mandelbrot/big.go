package mandelbrot

import (
	"image"
	"mandelbrot/fractal"
	"mandelbrot/palette"
	"math/big"
	"sync"
	"sync/atomic"
)

// Big is a mandelbrot fractal generator that uses *big.Float numbers for calculations
type Big struct {
	iterations int
	threshold  float32
}

func NewBigDefault() *Big {
	return NewBig(DefaultIterations, DefaultThreshold)
}

func NewBig(iterations int, threshold float32) *Big {
	ret := &Big{
		iterations: iterations,
		threshold:  threshold,
	}

	return ret
}

var two = big.NewFloat(2.0)
var half = big.NewFloat(0.5)

// Generation function
func (f *Big) Generate(
	target *image.RGBA,
	cx, cy, scale *big.Float,
	physicalWidth, physicalHeight *big.Float,
	reportingFunc fractal.ProgressReportingFunc,
	doneFunc fractal.DoneFunc,
) {
	go func() {
		// Start physical x point
		// physMinX = cx - (physWidth / 2)
		physMinX := big.NewFloat(0).Copy(physicalWidth)
		physMinX = physMinX.Mul(physMinX, half)
		physMinX = physMinX.Sub(cx, physMinX)

		// Start physical y point
		// physMinY = cy - (physHeight / 2)
		physMinY := big.NewFloat(0).Copy(physicalHeight)
		physMinY = physMinY.Mul(physMinY, half)
		physMinY = physMinY.Sub(cy, physMinY)

		// Calculate pixel-to-physical scale
		scaleX := big.NewFloat(0).Quo(physicalWidth, big.NewFloat(float64(target.Rect.Max.X)))
		scaleY := big.NewFloat(0).Quo(physicalHeight, big.NewFloat(float64(target.Rect.Max.Y)))

		pal := palette.CreatePaletteGrayscaleRecursive(256)

		wg := sync.WaitGroup{}

		// Counter will hold number of completed lines. Need for progress reporting
		var linesDone int32
		linesDonePtr := &linesDone

		// (x, y) - are pixel coords
		for y := 0; y < target.Rect.Max.Y; y++ {
			// (physX, physY) - are physical coordinates
			physY := big.NewFloat(float64(y)).SetPrec(cy.Prec())
			physY = physY.Mul(physY, scaleY)
			physY.Add(physY, physMinY)

			wg.Add(1)
			go func(y int, physY *big.Float) {
				for x := 0; x < target.Rect.Max.X; x++ {
					physX := big.NewFloat(float64(x)).SetPrec(cx.Prec())
					physX = physX.Mul(physX, scaleX)
					physX.Add(physX, physMinX)

					// get fractal value at the point
					value := mandelbrotBig(physX, physY, f.iterations, f.threshold)

					// convert it to the color and set pixel color
					target.Set(x, y, pal[int(float32(len(pal))*value)])
				}
				atomic.AddInt32(linesDonePtr, 1)
				reportingFunc(float32(atomic.LoadInt32(linesDonePtr)) / float32(target.Rect.Max.Y))
				wg.Done()
			}(y, physY)
		}

		wg.Wait()
		doneFunc()
	}()
}

func mandelbrotBig(x *big.Float, y *big.Float,iterations int, threshold float32) float32 {
	thresholdBig := big.NewFloat(float64(threshold))


	retX := big.NewFloat(0)
	retY := big.NewFloat(0)

	tmp := big.NewFloat(0)

	for i := 0; i < iterations; i++ {
		// calc real part: x^2 - y^2
		xSquared := big.NewFloat(0).Mul(retX, retX)
		ySquared := big.NewFloat(0).Mul(retY, retY)
		newRetX := big.NewFloat(0).Sub(xSquared, ySquared)

		// calc imaginary part: 2*x*y
		newRetY := big.NewFloat(0).Mul(retX, retY)
		newRetY.Mul(newRetY, two)

		// add (x, y)
		retX = newRetX.Add(newRetX, x)
		retY = newRetY.Add(newRetY, y)

		// calculate absolute value of complex number (retX, retY)
		// reuse previous calculations of xSquared and ySquared
		tmp.Add(xSquared, ySquared) // tmp <- x^2 + y^2
		tmp.Sqrt(xSquared) // tmp <- sqrt(tmp)

		abs := tmp

		if abs.Cmp(thresholdBig) > 0 {
			return float32(i) / float32(iterations)
		}
	}

	return 0
}
