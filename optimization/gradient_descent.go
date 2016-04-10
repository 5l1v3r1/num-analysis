package optimization

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
)

// GradFunc is a multivariable function with a
// computable gradient.
type GradFunc interface {
	// Dim returns the number of input arguments
	// this function takes.
	Dim() int

	// Eval evaluates the function for the given
	// argument vector.
	Eval(vec linalg.Vector) float64

	// Gradient evaluates the gradient of the
	// function at a given argument vector.
	Gradient(vec linalg.Vector) linalg.Vector
}

// GradientDescentPrec finds the minimum of a
// quasi-convex multivariable function using
// gradient descent.
// If a step is taken which moves the argument
// vector by less than prec (as measured by
// Euclidean distance), then the algorithm
// terminates.
func GradientDescent(f GradFunc, prec float64) linalg.Vector {
	guess := make(linalg.Vector, f.Dim())
	lastValue := f.Eval(guess)
	for {
		gradient := f.Gradient(guess)
		if vectorIsZero(gradient) {
			break
		}

		sf := stepSizeFunc{f, gradient, guess}
		stepSize := UnimodalMin(&sf)

		gradient.Scale(stepSize)
		guess.Add(gradient)

		value := f.Eval(guess)
		if value >= lastValue {
			// If the value isn't any lower, then we probably
			// can't get much better at machine precision.
			break
		}
		lastValue = value

		dist := math.Sqrt(gradient.Dot(gradient))
		if dist < prec {
			break
		}
	}
	return guess
}

type stepSizeFunc struct {
	f        GradFunc
	gradient linalg.Vector
	start    linalg.Vector
}

func (s *stepSizeFunc) Eval(x float64) float64 {
	return s.f.Eval(s.gradient.Copy().Scale(x).Add(s.start))
}

func vectorIsZero(v linalg.Vector) bool {
	for _, x := range v {
		if x != 0 {
			return false
		}
	}
	return true
}
