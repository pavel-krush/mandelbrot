package fractal

import (
	"fmt"
	"math"
	"math/big"
	"math/cmplx"
	"testing"
)

func TestBigLength(t *testing.T) {
	r := math.Pi
	i := math.E

	cplx := complex(r, i)
	cplxRes := cmplx.Abs(cplx)

	bigR := big.NewFloat(r)
	bigI := big.NewFloat(i)

	bigRes := bigLength(bigR, bigI)
	fmt.Printf("cplx res: %.20f\n", cplxRes)
	fmt.Printf("big res : %.20f\n", bigRes)
}
