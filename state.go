package main

import (
	"fmt"
	"math/big"
)

type State struct {
	precision                     uint       // Floats precision
	cx, cy                        *big.Float // Center point
	scale                         *big.Float // Scale ratio
	screenWidth, screenHeight     float64    // Screen width and height in pixels
	physicalWidth, physicalHeight *big.Float // Physical coordinates that is currently rendered
}

const (
	DefaultFloatsPrecision         = 80
	DefaultCenterX         float64 = -0.7
	DefaultCenterY         float64 = 0.0
	DefaultScale           float64 = 1.0
	DefaultScreenWidth     float64 = 640
	DefaultScreenHeight    float64 = 480
	DefaultPhysicalWidth   float64 = 3.0
	DefaultPhysicalHeight  float64 = 2.0
)

func NewState() *State {
	ret := &State{
		precision:      DefaultFloatsPrecision,
		cx:             big.NewFloat(DefaultCenterX),
		cy:             big.NewFloat(DefaultCenterY),
		scale:          big.NewFloat(DefaultScale),
		screenWidth:    DefaultScreenWidth,
		screenHeight:   DefaultScreenHeight,
		physicalWidth:  big.NewFloat(DefaultPhysicalWidth),
		physicalHeight: big.NewFloat(DefaultPhysicalHeight),
	}

	ret.SetPrecision(ret.precision)

	return ret
}

// Copy current state
func (s *State) Copy() *State {
	ret := &State{}

	ret.precision = s.precision
	ret.cx = big.NewFloat(0).Copy(s.cx)
	ret.cy = big.NewFloat(0).Copy(s.cy)
	ret.scale = big.NewFloat(0).Copy(s.scale)
	ret.screenWidth, ret.screenHeight = s.screenWidth, s.screenHeight
	ret.physicalWidth = big.NewFloat(0).Copy(s.physicalWidth)
	ret.physicalHeight = big.NewFloat(0).Copy(s.physicalHeight)

	return ret
}

func (s *State) String() string {
	return "State(\n" +
		fmt.Sprintf(" cx=%.50f\n", s.cx) +
		fmt.Sprintf(" cy=%.50f\n", s.cy) +
		fmt.Sprintf(" scale=%.50f\n", s.scale) +
		fmt.Sprintf(" physWidth=%.50f\n", s.physicalWidth) +
		fmt.Sprintf(" physHeight=%.50f\n", s.physicalHeight) +
		")"
}

// Change precision
func (s *State) SetPrecision(precision uint) {
	s.precision = precision
	s.cx.SetPrec(s.precision)
	s.cy.SetPrec(s.precision)
	s.scale.SetPrec(s.precision)
}

// Get current floats precision
func (s *State) GetPrecision() uint {
	return s.precision
}

// Return current center x-coordinate
func (s *State) GetCX() *big.Float {
	return s.cx
}

// Return current center y-coordinate
func (s *State) GetCY() *big.Float {
	return s.cy
}

// Return current scale
func (s *State) GetScale() *big.Float {
	return s.scale
}

// Return rendering screen width
func (s *State) GetScreenWidth() float64 {
	return s.screenWidth
}

// Return rendering screen width
func (s *State) GetScreenHeight() float64 {
	return s.screenHeight
}

// Return rendering physical width
func (s *State) GetPhysicalWidth() *big.Float {
	return s.physicalWidth
}

// Return rendering physical height
func (s *State) GetPhysicalHeight() *big.Float {
	return s.physicalHeight
}
