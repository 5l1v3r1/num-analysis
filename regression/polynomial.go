package regression

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/cholesky"
)

// Polynomial represents a polynomial with real coefficients.
// The first entry is the coefficient to x^0, then x^1, etc.
type Polynomial []float64

// FitPolynomial uses least-squares to find a polynomial
// to match the set of points as well as possible.
//
// If there are fewer points than deg?1, or if there are
// fewer than deg+1 unique Input values, then this may
// return an incorrect or invalid solution.
func FitPolynomial(deg int, points []Point) Polynomial {
	outputs := make(linalg.Vector, len(points))
	for row, p := range points {
		outputs[row] = p.Output
	}

	degInputs := make([]linalg.Vector, deg+1)
	for d := range degInputs {
		v := make(linalg.Vector, len(points))
		for i, p := range points {
			v[i] = math.Pow(p.Input, float64(i))
		}
		degInputs[d] = v
	}

	normalMat := linalg.NewMatrix(deg+1, deg+1)
	for row, v1 := range degInputs {
		for col, v2 := range degInputs {
			normalMat.Set(row, col, v1.Dot(v2))
		}
	}

	normalOut := make(linalg.Vector, deg+1)
	for i := range normalOut {
		normalOut[i] = degInputs[i].Dot(outputs)
	}

	return Polynomial(cholesky.Decompose(normalMat).Solve(normalOut))
}
