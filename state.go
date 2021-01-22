package main

import "math/big"

type State struct {
	precision uint       // floats precision
	cx, cy    *big.Float // center point
	scale     *big.Float // zoom level
}

const (
	DefaultFloatsPrecision         = 8
	DefaultCenterX         float64 = -0.7
	DefaultCenterY         float64 = 0.0
	DefaultScale           float64 = 1.0
)

func NewState() *State {
	ret := &State{
		precision: DefaultFloatsPrecision,
		cx:        big.NewFloat(DefaultCenterX),
		cy:        big.NewFloat(DefaultCenterY),
		scale:     big.NewFloat(DefaultScale),
	}

	ret.SetPrecision(ret.precision)

	return ret
}

// Change precision
func (s *State) SetPrecision(precision uint) {
	s.precision = precision
	s.cx.SetPrec(s.precision)
	s.cy.SetPrec(s.precision)
	s.scale.SetPrec(s.precision)
}

// Set center coordinates and a scale factor
func (s *State) SetCoords(cx, cy, scale *big.Float) {
	s.cx = cx
	s.cy = cy
	s.scale = scale
}
