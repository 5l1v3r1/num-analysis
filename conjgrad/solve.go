package conjgrad

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
)

const residualUpdateFrequency = 20

// SolvePrec solves a system of linear equations
// t*x = b for x, where t is a symmetric
// positive-definite linear transformation.
//
// The prec argument specifies a bound on the
// residual error of the solution. If the largest
// element of (Ax-b) has an absolute value less than
// prec, then the current x is returned.
func SolvePrec(t LinTran, b linalg.Vector, prec float64) linalg.Vector {
	var conjVec linalg.Vector
	var residual linalg.Vector
	var solution linalg.Vector

	residual = b
	solution = make(linalg.Vector, t.Dim())

	for i := 0; i < t.Dim(); i++ {
		if greatestValue(residual) <= prec {
			break
		}
		if i == 0 {
			conjVec = residual.Copy()
		} else {
			// TODO: derive more efficient formula for projecting conjVec onto residual.
			projAmount := residual.Dot(t.Apply(conjVec)) / conjVec.Dot(t.Apply(conjVec))
			conjVec = residual.Copy().Add(conjVec.Scale(-projAmount))
		}
		if allZero(conjVec) {
			break
		}
		optimalDistance := conjVec.Dot(residual) / conjVec.Dot(t.Apply(conjVec))

		solution.Add(conjVec.Copy().Scale(optimalDistance))
		if i != 0 && (i%residualUpdateFrequency) == 0 {
			residual = t.Apply(solution).Scale(-1).Add(b)
		} else {
			residual.Add(t.Apply(conjVec).Scale(-optimalDistance))
		}
	}

	return solution
}

// Solve is like SolvePrec, except that it computes
// as accurate a solution as possible.
func Solve(t LinTran, b linalg.Vector) linalg.Vector {
	return SolvePrec(t, b, 0)
}

func allZero(v linalg.Vector) bool {
	for _, x := range v {
		if x != 0 {
			return false
		}
	}
	return true
}

func greatestValue(v linalg.Vector) float64 {
	var max float64
	for _, x := range v {
		max = math.Max(max, math.Abs(x))
	}
	return max
}
