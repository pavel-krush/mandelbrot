package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	x := big.NewFloat(math.Pi)
	squared := big.NewFloat(0).Mul(x, x)

	fmt.Printf("x          : %.10f\n", x)
	fmt.Printf("x ptr      : %p\n", x)
	fmt.Printf("squared    : %.10f\n", squared)
	fmt.Printf("squared ptr: %p\n", squared)
}
