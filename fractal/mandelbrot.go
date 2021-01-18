package fractal

import (
	"math/big"
	"math/cmplx"
)

func Mandelbrot(c complex128) float64 {
	const iterations = 256
	const threshold = 3.0

	ret := complex(0, 0)
	for i := 0; i < iterations; i++ {
		ret = ret*ret + c
		if cmplx.Abs(ret) > threshold {
			return float64(i) / float64(iterations)
		}
	}

	return 0
}

var two = big.NewFloat(2.0)

func MandelbrotBig(x *big.Float, y *big.Float) float64 {
	const iterations = 256
	const threshold = 3.0

	thresholdBig := big.NewFloat(threshold)

	retX := big.NewFloat(0)
	retY := big.NewFloat(0)

	for i := 0; i < iterations; i++ {
		// calc real part: x^2 - y^2
		xSquared := big.NewFloat(0).Mul(retX, retX)
		ySquared := big.NewFloat(0).Mul(retY, retY)
		newRetX := big.NewFloat(0).Sub(xSquared, ySquared)

		// calc imaginary part: 2*x*y
		xy := big.NewFloat(0).Mul(retX, retY)
		newRetY := big.NewFloat(0).Mul(xy, two)

		// add (x, y)
		retX = newRetX.Add(newRetX, x)
		retY = newRetY.Add(newRetY, y)

		// calculate absolute value of complex number (retX, retY)
		// reuse previous calculations of xSquared and ySquared
		sumOfSquares := big.NewFloat(0).Add(xSquared, ySquared)
		abs := sumOfSquares.Sqrt(sumOfSquares)

		if abs.Cmp(thresholdBig) > 0 {
			return float64(i) / float64(iterations)
		}
	}

	return 0
}
