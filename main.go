package main

import (
	"flag"
	"fmt"
	"mandelbrot/fractal/mandelbrot"
	"math/big"
)

func MustParseBigFloat(s string, precision uint) *big.Float {
	z, _, err := big.ParseFloat(s, 10, precision, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return z
}

func main() {
	app := NewApplication("Mandelbrot Fractal Explorer")

	generatorStr := flag.String("generator", "float64", "select generator: big or float64")
	flag.Parse()

	if generatorStr == nil {
		panic("generator")
	}

	if *generatorStr == "big" {
		app.SetGenerator(mandelbrot.NewBigDefault())
	} else if *generatorStr == "float64" {
		app.SetGenerator(mandelbrot.NewFloat64Default())
	} else {
		panic(*generatorStr)
	}

	fmt.Printf("Using %s generator\n", *generatorStr)

	//cx := "-1.48656573768883788853042260418005804552266102547264"
	//cy := "0.03579713550865033095370105522259793185378684565734"
	//scale := "0.00000000000000640180414098903887916742577864421037"
	//physWidth := "0.00000000000001920541242296711114025336924125917013"
	//physHeight := "0.00000000000001280360828197809092253489087575796220"
	//
	//precision := application.state.GetPrecision()
	//
	//application.state.GetCX().Set(MustParseBigFloat(cx, precision))
	//application.state.GetCY().Set(MustParseBigFloat(cy, precision))
	//application.state.GetScale().Set(MustParseBigFloat(scale, precision))
	//application.state.GetPhysicalWidth().Set(MustParseBigFloat(physWidth, precision))
	//application.state.GetPhysicalHeight().Set(MustParseBigFloat(physHeight, precision))

	app.Run()
}
