package big

import "math/big"

func MustParseBigFloat(s string, precision uint) *big.Float {
	z, _, err := big.ParseFloat(s, 10, precision, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return z
}
