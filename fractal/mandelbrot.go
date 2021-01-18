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

func bigLength(x *big.Float, y *big.Float) *big.Float {
	xSquared := big.NewFloat(0).Mul(x, x)
	ySquared := big.NewFloat(0).Mul(y, y)
	sumOfSquares := big.NewFloat(0).Add(xSquared, ySquared)
	return big.NewFloat(0).Sqrt(sumOfSquares)
}

func MandelbrotBig(x *big.Float, y *big.Float) float64 {
	const iterations = 256
	const threshold = 3.0

	thresholdBig := big.NewFloat(threshold)

	two := big.NewFloat(2.0)

	retX := big.NewFloat(0).SetPrec(x.Prec())
	retY := big.NewFloat(0).SetPrec(y.Prec())

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

		length := bigLength(retX, retY)

		if length.Cmp(thresholdBig) > 0 {
			return float64(i) / float64(iterations)
		}
	}

	return 0
}
