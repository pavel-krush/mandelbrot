package mandelbrot

import (
	"math/big"
	"testing"
)

func getTestingParams() (*big.Float, *big.Float, *big.Float, *big.Float, int, int) {
	const precision = 100

	cx := big.NewFloat(-0.7).SetPrec(precision)
	cy := big.NewFloat(0.0).SetPrec(precision)
	physWidth := big.NewFloat(3.0).SetPrec(precision)
	physHeight := big.NewFloat(2.0).SetPrec(precision)
	screenWidth := 640
	screenHeight := 480

	return cx, cy, physWidth, physHeight, screenWidth, screenHeight
}

func BenchmarkFloat64(b *testing.B) {
	cx, cy, physicalWidth, physicalHeight, screenWidth, screenHeight := getTestingParams()
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
	scaleX := big.NewFloat(0).SetPrec(physicalWidth.Prec()).Quo(physicalWidth, big.NewFloat(float64(screenWidth)))
	scaleY := big.NewFloat(0).SetPrec(physicalHeight.Prec()).Quo(physicalHeight, big.NewFloat(float64(screenHeight)))

	points := make([]struct{x, y float64}, screenHeight * screenWidth)

	// (x, y) - are pixel coords
	for y := 0; y < screenHeight; y++ {
		// (physX, physY) - are physical coordinates
		physY := big.NewFloat(float64(y)).SetPrec(cy.Prec())
		physY = physY.Mul(physY, scaleY)
		physY.Add(physY, physMinY)
		for x := 0; x < screenWidth; x++ {
			physX := big.NewFloat(float64(x)).SetPrec(cx.Prec())
			physX = physX.Mul(physX, scaleX)
			physX.Add(physX, physMinX)

			points[y * screenWidth + x].x, _ = physX.Float64()
			points[y * screenWidth + x].y, _ = physY.Float64()
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		index := i % len(points)
		mandelbrotComplex128(complex(points[index].x, points[index].y), 10, 3.0)
	}
	b.StopTimer()
}

func BenchmarkBig(b *testing.B) {
	cx, cy, physicalWidth, physicalHeight, screenWidth, screenHeight := getTestingParams()
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
	scaleX := big.NewFloat(0).SetPrec(physicalWidth.Prec()).Quo(physicalWidth, big.NewFloat(float64(screenWidth)))
	scaleY := big.NewFloat(0).SetPrec(physicalHeight.Prec()).Quo(physicalHeight, big.NewFloat(float64(screenHeight)))

	points := make([]struct{x, y *big.Float}, screenHeight * screenWidth)

	// (x, y) - are pixel coords
	for y := 0; y < screenHeight; y++ {
		// (physX, physY) - are physical coordinates
		physY := big.NewFloat(float64(y)).SetPrec(cy.Prec())
		physY = physY.Mul(physY, scaleY)
		physY.Add(physY, physMinY)
		for x := 0; x < screenWidth; x++ {
			physX := big.NewFloat(float64(x)).SetPrec(cx.Prec())
			physX = physX.Mul(physX, scaleX)
			physX.Add(physX, physMinX)

			points[y * screenWidth + x].x = physX
			points[y * screenWidth + x].y = physY
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		index := i % len(points)
		mandelbrotBig(points[index].x, points[index].y, 10, 3.0)
	}
	b.StopTimer()
}
