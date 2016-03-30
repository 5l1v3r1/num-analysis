package ludecomp

import (
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestBackSubstitution(t *testing.T) {
	m1 := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			3, 5, 7, 2,
			1, 9.5, 2, 4,
			3, 18, 3.6, 5,
			2, 8, 1, 3.14,
		},
	}
	lowerSolution := solveLowerTriangular(m1, linalg.Vector{1, 2, 3, 4})
	upperSolution := solveUpperTriangular(m1, linalg.Vector{1, 2, 3, 4})

	expectedLower := linalg.Vector{1.0 / 3.0, 10.0 / 57.0, -55.0 / 171.0, 6317.0 / 8810.0}
	expectedUpper := linalg.Vector{12, 20, -17, 4}
	if vectorDiff(lowerSolution, expectedLower) > 0.0001 {
		t.Error("incorrect lower triangular solution", lowerSolution)
	}
	if vectorDiff(upperSolution, expectedUpper) > 0.0001 {
		t.Error("incorrect upper triangular solution", upperSolution)
	}
}

func vectorDiff(v1, v2 linalg.Vector) float64 {
	return v2.Copy().Add(v1.Copy().Scale(-1)).Mag()
}
