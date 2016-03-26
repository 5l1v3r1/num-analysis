package ludecomp

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/kahan"
)

func TestBackSubstitution(t *testing.T) {
	m1 := Matrix{
		N: 4,
		V: []float64{
			3, 5, 7, 2,
			1, 9.5, 2, 4,
			3, 18, 3.6, 5,
			2, 8, 1, 3.14,
		},
	}
	lowerSolution := m1.SolveLowerTriangular(Vector{1, 2, 3, 4})
	upperSolution := m1.SolveUpperTriangular(Vector{1, 2, 3, 4})

	expectedLower := Vector{1.0 / 3.0, 10.0 / 57.0, -55.0 / 171.0, 6317.0 / 8810.0}
	expectedUpper := Vector{12, 20, -17, 4}
	if vectorDiff(lowerSolution, expectedLower) > 0.0001 {
		t.Error("incorrect lower triangular solution", lowerSolution)
	}
	if vectorDiff(upperSolution, expectedUpper) > 0.0001 {
		t.Error("incorrect upper triangular solution", upperSolution)
	}
}

func vectorDiff(v1, v2 Vector) float64 {
	diffList := make([]float64, len(v1))
	for i, x := range v1 {
		diffList[i] = math.Pow(x-v2[i], 2)
	}
	return math.Sqrt(kahan.Sum64(diffList))
}
