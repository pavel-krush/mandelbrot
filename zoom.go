package main

import (
	"math/big"
)

type ZoomDirection int

const (
	ZoomDirectionIn ZoomDirection = iota
	ZoomDirectionOut
)

type Zoomer interface {
	// ZoomAt should modify given state to make it zoomed at point x, y
	ZoomAt(s *State, x, y float64, zoomDirection ZoomDirection)
}

type ZoomerSimple struct {}

func NewZoomerSimple() *ZoomerSimple {
	return &ZoomerSimple{}
}

// Calculate new center coordinates from the clicked coordinates
func (z *ZoomerSimple) ZoomAt(
	s *State, x, y float64, zoomDirection ZoomDirection,
) {
	const zoomFactorDelta = 0.1

	// Apply only that factor of the real offset to make zooming smoother
	offsetScale := big.NewFloat(0.2)

	// Rescale coordinates from [0, screenWidth] to [-1, 1]
	normX := ((x / s.GetScreenWidth()) - 0.5) * 2
	normY := ((y / s.GetScreenHeight()) - 0.5) * 2

	// Calculate new center coordinates: x
	cx, cy, scale := s.GetCX(), s.GetCY(), s.GetScale()
	newX := big.NewFloat(normX).SetPrec(cx.Prec())
	newX.Mul(newX, s.GetPhysicalWidth())
	newX.Mul(newX, offsetScale)
	newX.Add(newX, cx)
	cx.Copy(newX)

	// Calculate new center coordinates: y
	newY := big.NewFloat(normY).SetPrec(cy.Prec())
	newY.Mul(newY, s.GetPhysicalHeight())
	newY.Mul(newY, offsetScale)
	newY.Add(newY, cy)
	cy.Copy(newY)

	// Adjust scale factor
	var zoomFactor float64
	switch zoomDirection {
	case ZoomDirectionIn:
		zoomFactor = 1 - zoomFactorDelta
	case ZoomDirectionOut:
		zoomFactor = 1 + zoomFactorDelta
	default:
		panic(zoomDirection)
	}

	// Adjust scale
	zoomFactorBig := big.NewFloat(zoomFactor).SetPrec(scale.Prec())
	scale.Mul(scale, zoomFactorBig)

	// Recalculate physical bounds
	physWidth := s.GetPhysicalWidth()
	physWidth.Mul(physWidth, big.NewFloat(zoomFactor))

	physHeight := s.GetPhysicalHeight()
	physHeight.Mul(physHeight, big.NewFloat(zoomFactor))
}
