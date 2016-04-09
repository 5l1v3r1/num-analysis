package mvroots

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/ludecomp"
)

// Iterator uses a multivariable extension of
// Newton's method to search for the root of
// a differentiable multivariable function.
type Iterator struct {
	function Func
	guess    linalg.Vector
}

// NewIterator creates an Iterator with a
// given input vector at which to start the
// search.
func NewIterator(f Func, start linalg.Vector) *Iterator {
	return &Iterator{function: f, guess: start}
}

// Step performs one root-finding iteration.
// It returns the Euclidean distance between
// the previous guess and the current guess.
func (i *Iterator) Step() float64 {
	value := i.function.Eval(i.guess)
	zero := true
	for _, x := range value {
		if x != 0 {
			zero = false
			break
		}
	}
	if zero {
		return 0
	}

	jacobian := i.function.Jacobian(i.guess)
	lu := ludecomp.Decompose(jacobian)

	// If the Jacobian is not invertible, then
	// there isn't necessarily a way to get to
	// zero from here.
	if lu.PivotScale() < math.Nextafter(1, 2)-1 {
		return 0
	}

	diff := lu.Solve(value).Scale(-1)
	diffMag := math.Sqrt(diff.Dot(diff))
	i.guess = diff.Add(i.guess)
	return diffMag
}

// Guess returns the current approximate root.
func (i *Iterator) Guess() linalg.Vector {
	return i.guess
}
