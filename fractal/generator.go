package fractal

import (
	"image"
	"math/big"
)

type ProgressReportingFunc func(progress float32)
type DoneFunc func()

// FractalGenerator defines an interface for objects that can draw fractals
type Generator interface {
	// target - fractal will be rendered here
	// cx, cy, scale - center coordinates and scale
	// reportingFunc - callback that could be called during generation
	// doneFunc - callback that must be called once after the generation is complete
	Generate(
		target *image.RGBA,
		cx, cy, scale *big.Float,
		physicalWidth, physicalHeight *big.Float,
		reportingFunc ProgressReportingFunc, doneFunc DoneFunc,
	)
}
