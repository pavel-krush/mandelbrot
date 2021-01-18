package main

import (
	"fmt"
	"math/big"
)

func mustParseBigFloat(s string) *big.Float {
	const precision = 1024
	z, _, err := big.ParseFloat(s, 10, precision, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return z
}

func main() {
	cxStr := "-0.74920101504"
	cyStr := "-0.0999999899"
	scaleStr := "1.0"
	zoomFactorStr := "0.95"
	minScaleStr := "0.0000000000001"

	cx := mustParseBigFloat(cxStr)
	cy := mustParseBigFloat(cyStr)
	scale := mustParseBigFloat(scaleStr)
	zoomFactor := mustParseBigFloat(zoomFactorStr)
	minScale := mustParseBigFloat(minScaleStr)

	i := 0
	for scale.Cmp(minScale) > 0 {
		filename := fmt.Sprintf("images/simple-%03d.png", i)
		fmt.Printf("go run simple.go %.100f %.100f %.100f > %s\n", cx, cy, scale, filename)
		i++
		scale.Mul(scale, zoomFactor)
		fmt.Printf("%d - %.100f\n", i, scale)
	}
}
