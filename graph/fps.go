package graph

import (
	"math"
	"time"
)

const (
	FPSMax             = 10000.0
	FPSSmoothingFactor = 0.00
)

type FPS struct {
	lastTick time.Time
	avg      float64
}

func (fps *FPS) FrameRendered() {
	currentTick := time.Now()

	currentFPS := 1 / math.Max(float64(currentTick.Sub(fps.lastTick))/float64(time.Second), 1/FPSMax)

	fps.lastTick = currentTick

	fps.avg = (fps.avg * FPSSmoothingFactor) + (currentFPS * (1 - FPSSmoothingFactor))
}

func (fps *FPS) GetFPS() float64 {
	return fps.avg
}
