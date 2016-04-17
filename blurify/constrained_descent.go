package main

import (
	"math"

	"github.com/unixpickle/num-analysis/conjgrad"
	"github.com/unixpickle/num-analysis/linalg"
)

// Descender uses an active set optimization method
// to find an image which is as close to the
// unblurred image as possible.
type Descender struct {
	blurred   linalg.Vector
	blurOp    conjgrad.LinTran
	bestGuess linalg.Vector

	temp []int
}

func NewDescender(blurred linalg.Vector, op conjgrad.LinTran) *Descender {
	res := &Descender{
		blurred:   blurred,
		blurOp:    op,
		bestGuess: make(linalg.Vector, len(blurred)),
		temp:      make([]int, len(blurred)),
	}
	for i := range res.bestGuess {
		res.bestGuess[i] = 0.5
	}
	return res
}

// Step performs one iteration of the algorithm,
// returning the sum of the differences between
// the old guess and the new one.
func (d *Descender) Step() float64 {
	constraints := d.constraints()
	residual := d.blurOp.Apply(d.bestGuess).Scale(-1).Add(d.blurred)
	allowedDirection := residual.Copy()
	for i, x := range constraints {
		if x == -1 {
			if residual[i] < 0 {
				allowedDirection[i] = 0
			}
		} else if x == 1 {
			if residual[i] > 0 {
				allowedDirection[i] = 0
			}
		}
	}

	optDist := residual.Dot(allowedDirection) /
		d.blurOp.Apply(allowedDirection).Dot(allowedDirection)

	oldGuess := d.bestGuess.Copy()
	d.bestGuess.Add(allowedDirection.Scale(optDist))

	var diffSum float64
	for i, x := range d.bestGuess {
		d.bestGuess[i] = math.Max(0, math.Min(1, x))
		diffSum += math.Abs(d.bestGuess[i] - oldGuess[i])
	}
	return diffSum
}

// Guess returns the current approximation of the
// unblurred image.
func (d *Descender) Guess() linalg.Vector {
	return d.bestGuess
}

func (d *Descender) constraints() []int {
	for i, x := range d.bestGuess {
		if x == 0 {
			d.temp[i] = -1
		} else if x == 1 {
			d.temp[i] = 1
		} else {
			d.temp[i] = 0
		}
	}
	return d.temp
}
